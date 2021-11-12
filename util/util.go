package util

import (
	"archive/zip"
	"context"
	"crypto/tls"
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"math"
	"math/rand"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"

	"google.golang.org/api/oauth2/v2"
	"google.golang.org/api/option"
	"streetbox.id/cfg"
	"streetbox.id/model"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	"golang.org/x/crypto/acme/autocert"
	"golang.org/x/crypto/bcrypt"
)

const (
	// TasksRegular enum
	TasksRegular = "REGULAR"
	// TasksNonRegular enum
	TasksNonRegular = "NONREGULAR"
	// TasksHomeVisit enum
	TasksHomeVisit = "HOMEVISIT"
	// TrxOrder enum
	TrxOrder = "ORDER"
	// TrxHomeVisit enum
	TrxHomeVisit = "VISIT"
	// TrxStatusSuccess success
	TrxStatusSuccess = "SUCCESS"
	// TrxStatusPending ..
	TrxStatusPending = "PENDING"
	// TrxStatusFailed ..
	TrxStatusFailed = "FAILED"
	// TrxStatusVoid TBA
	TrxStatusVoid = "VOID"
	// TrxOrderOnlineTypes ..
	TrxOrderOnlineTypes = 1
	// TrxOrderOfflineTypes ..
	TrxOrderOfflineTypes = 0
	// TrxRefundParkingSpace ..
	TrxRefundParkingSpace = "SPACE"
	// TrxRefundHomeVisit ..
	TrxRefundHomeVisit = "VISIT"
	// TrxVisitStatusOpen enum
	TrxVisitStatusOpen = "OPEN"
	// TrxVisitStatusClosed enum
	TrxVisitStatusClosed = "CLOSED"
	// RoleSuper enum
	RoleSuper = "superadmin"
	// RoleAdmin enum
	RoleAdmin = "admin"
	// RoleMerchant enum
	RoleMerchant = "merchant"
	// RoleFoodtruck enum
	RoleFoodtruck = "foodtruck"
	// RoleConsumer enum
	RoleConsumer = "consumer"
)

// HashPassword ...
func HashPassword(password string) string {
	bytes, _ := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes)
}

// CheckPasswordHash ...
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// CreateToken ...
func CreateToken(id int64, roleName string) string {
	if id == 0 || roleName == "" {
		return ""
	}
	jwtKey := cfg.Config.JwtKey
	claims := jwt.MapClaims{}
	claims["user_id"] = id
	claims["role_name"] = roleName
	// if roleName == "consumer" || roleName == "foodtruck" {
	// 	claims["exp"] = time.Now().Add(time.Hour * 8766).Unix()
	// }
	// claims["exp"] = time.Now().Add(time.Hour * 1).Unix()

	// token never expired
	claims["exp"] = 0
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(jwtKey))
	if err != nil {
		return ""
	}
	return signedToken
}

// CreateTokenReset ...
func CreateTokenReset(id int64, roleName string) string {
	jwtKey := cfg.Config.JwtKey
	claims := jwt.MapClaims{}
	claims["user_id"] = id
	claims["role_name"] = roleName
	claims["exp"] = time.Now().Add(time.Minute * 30).Unix()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(jwtKey))
	if err != nil {
		return ""
	}
	return signedToken
}

// ExtractToken get token from header
func ExtractToken(c *http.Request) string {
	bearToken := c.Header.Get("Authorization")
	strArr := strings.Split(bearToken, " ")
	if len(strArr) == 2 {
		return strArr[1]
	}
	return ""
}

func verifyToken(c *http.Request) (*jwt.Token, error) {
	jwtKey := cfg.Config.JwtKey
	tokenString := ExtractToken(c)
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(jwtKey), nil
	})
	if err != nil {
		return nil, err
	}
	return token, nil
}

// TokenValid token and permission role validation
//
// 1. all		: all authorized user
//
// 2. superadmin: superadmin
//
// 3. merchant	: foodtruck -> admin -> superadmin
//
// 4. admin		: admin -> superadmin
//
// 5. consumer	: consumer
func TokenValid(c *http.Request, permission string) error {
	token, err := verifyToken(c)
	if err != nil {
		return errors.New("Unauthorized Access, Please Login First")
	}
	if _, ok := token.Claims.(jwt.Claims); !ok && !token.Valid {
		return errors.New("Unauthorized Access, Please Login First")
	}
	claims := token.Claims.(jwt.MapClaims)
	currentRole := claims["role_name"].(string)
	if permission == "all" {
		return nil
	}
	if permission == RoleSuper {
		if currentRole == RoleSuper {
			return nil
		}
	} else if permission == RoleMerchant {
		if currentRole == RoleFoodtruck ||
			currentRole == RoleAdmin ||
			currentRole == RoleSuper {
			return nil
		}
	} else if permission == RoleAdmin {
		if currentRole == RoleAdmin ||
			currentRole == RoleSuper {
			return nil
		}
	} else if permission == RoleConsumer {
		if currentRole == RoleConsumer {
			return nil
		}
	}
	return errors.New("Unauthorized Role Access")
}

// ExtractTokenMetadata ...
func ExtractTokenMetadata(c *http.Request) (*model.JwtCustomClaims, error) {
	token, err := verifyToken(c)
	if err != nil {
		return nil, err
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if ok && token.Valid {
		roleName, ok := claims["role_name"].(string)
		if !ok {
			return nil, err
		}
		userID, err := strconv.ParseInt(fmt.Sprintf("%.f", claims["user_id"]), 10, 64)
		if err != nil {
			return nil, err
		}
		return &model.JwtCustomClaims{
			UserID:   userID,
			RoleName: roleName,
		}, nil
	}
	return nil, err
}

// ParamToDatetime ...
func ParamToDatetime(param string) time.Time {
	layout := "2006-01-02 15:04:05"
	t, err := time.Parse(layout, param)
	if err != nil {
		log.Panicln(err)
	}
	return t
}

// ParamToDate ...
func ParamToDate(param string) time.Time {
	layout := "2006-01-02"
	t, err := time.Parse(layout, param)
	if err != nil {
		log.Panicln(err)
	}
	return t
}

// ParamIDToInt64 ...
func ParamIDToInt64(param string) int64 {
	result, err := strconv.ParseInt(param, 10, 64)
	if err != nil {
		log.Panicln(err)
		return 0
	}
	return result
}

// ParamToFloat64 ...
func ParamToFloat64(param string) float64 {
	result, err := strconv.ParseFloat(param, 64)
	if err != nil {
		log.Panicln(err)
		return 0
	}
	return result
}

// ParamToFloat32 ...
func ParamToFloat32(param string) float32 {
	result, err := strconv.ParseFloat(param, 32)
	if err != nil {
		log.Panicln(err)
		return 0
	}
	return float32(result)
}

// ParamIDToInt ...
func ParamIDToInt(param string) int {
	result, err := strconv.Atoi(param)
	if err != nil {
		log.Panicln(err)
		return 0
	}
	return result
}

// GeneratedUUID renaming uploaded name files to UUID
func GeneratedUUID(filename string) string {
	names := strings.Split(filename, ".")
	uuid := uuid.New().String()
	return fmt.Sprintf("%s.%s", uuid, names[len(names)-1])
}

// SortedBy ..
func SortedBy(sort []string) []string {
	var sorted []string
	for _, v := range sort {
		split := strings.Split(v, ",")
		sorted = append(sorted, fmt.Sprintf("%s %s", split[0], split[1]))
	}
	return sorted
}

// Offset ..
func Offset(page, limit int) int {
	offset := 0
	if page == 1 {
		offset = 0
	} else {
		offset = (page - 1) * limit
	}
	return offset
}

// TotalPages ..
func TotalPages(count, limit int) int {
	return int(math.Ceil(float64(count) / float64(limit)))
}

// NextPage ..
func NextPage(page, totalPages int) int {
	if page == totalPages {
		return page
	}
	return page + 1
}

// PrevPage ..
func PrevPage(page int) int {
	if page > 1 {
		return page - 1
	}
	return page
}

// DateTimeToMilliSeconds ...
func DateTimeToMilliSeconds(date time.Time) int64 {
	return date.UnixNano() / int64(time.Millisecond)
}

// DateStringToMilliSeconds ...
func DateStringToMilliSeconds(date string) int64 {
	dates, _ := time.Parse("2006-01-02 15:04:05", date)
	return DateTimeToMilliSeconds(dates)
}

// DateTimeSwap ..
func DateTimeSwap(date1 time.Time, date2 time.Time) string {
	datesplit1 := strings.Split(date1.Format("2006-01-02 15:04:05"), " ")[0]
	datesplit2 := strings.Split(date2.Format("2006-01-02 15:04:05"), " ")[1]
	return datesplit1 + " " + datesplit2
}

// ExtractGoogleTokenInfo ..
func ExtractGoogleTokenInfo(idToken string) (*oauth2.Tokeninfo, error) {
	optClient := option.WithHTTPClient(&http.Client{})
	oauthSvc, err := oauth2.NewService(context.Background(), optClient)
	if err != nil {
		return nil, err
	}
	tokenInfoCall := oauthSvc.Tokeninfo()
	tokenInfoCall.IdToken(idToken)
	return tokenInfoCall.Do()
}

// GetSelfSignedOrLetsEncryptCert ...
func GetSelfSignedOrLetsEncryptCert(
	certMan *autocert.Manager) func(hello *tls.ClientHelloInfo) (*tls.Certificate, error) {
	return func(hello *tls.ClientHelloInfo) (*tls.Certificate, error) {
		dirCache, ok := certMan.Cache.(autocert.DirCache)
		certsDir := cfg.Config.Certificate.Dir
		certFileName := cfg.Config.Certificate.Filename
		if !ok {
			dirCache = autocert.DirCache(certsDir)
		}
		keyFile := filepath.Join(string(dirCache), certFileName+".key")
		crtFile := filepath.Join(string(dirCache), certFileName+".crt")
		certificate, err := tls.LoadX509KeyPair(crtFile, keyFile)
		if err != nil {
			fmt.Printf("%s\nFalling back to Letsencrypt\n", err)
			return certMan.GetCertificate(hello)
		}
		fmt.Println("Loaded selfsigned certificate.")
		return &certificate, err
	}
}

// GenerateTrxID -> TRX-XXXX9999
func GenerateTrxID() string {
	const charStr = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	const charNum = "0123456789"
	var seededRand *rand.Rand = rand.New(rand.NewSource(time.Now().UnixNano()))
	rndCharStr := make([]byte, 4)
	rndCharNum := make([]byte, 4)
	for r := range rndCharStr {
		rndCharStr[r] = charStr[seededRand.Intn(len(charStr))]
	}
	for r := range rndCharNum {
		rndCharNum[r] = charNum[seededRand.Intn(len(charNum))]
	}
	return "TRX-" + string(rndCharStr) + string(rndCharNum)
}

// DoRequest ...
func DoRequest(req *http.Request) ([]byte, error) {
	timeout := time.Duration(20 * time.Second)
	client := http.Client{
		Timeout: timeout,
	}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("ERROR: %s", err.Error())
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("ERROR: %s", err.Error())
		return nil, err
	}
	return body, nil
}

// MillisToTime ..
func MillisToTime(millis int64) time.Time {
	return time.Unix(0, millis*int64(time.Millisecond))
}

// Unzip ..
// Unzip will decompress a zip archive, moving all files and folders
// within the zip file (parameter 1) to an output directory (parameter 2).
func Unzip(src *multipart.FileHeader, dest string, c *gin.Context) ([]string, []string, error) {

	var (
		filenames     []string
		uuidFileNames []string
	)

	// check if os windows
	if runtime.GOOS == "windows" {
		drive := "D:"
		replacedest := strings.ReplaceAll(dest, "/", "\\")
		dest = drive + replacedest
	}

	fileName := src.Filename
	path := dest + fileName

	if err := c.SaveUploadedFile(src, path); err != nil {
		return nil, nil, err
	}

	r, err := zip.OpenReader(path)
	if err != nil {
		return nil, nil, err
	}
	defer r.Close()

	for _, f := range r.File {

		// Store filename/path for returning and using later on
		filename := GeneratedUUID(filepath.Base(f.Name))
		fpath := filepath.Join(dest, filename)

		// Check for ZipSlip. More Info: http://bit.ly/2MsjAWE
		if !strings.HasPrefix(fpath, filepath.Clean(dest)+string(os.PathSeparator)) {
			return nil, nil, fmt.Errorf("%s: illegal file path", fpath)
		}

		outFile, err := os.OpenFile(fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			fmt.Println("error : os open file")
			return nil, nil, err
		}
		rc, err := f.Open()
		if err != nil {
			fmt.Println("error : f open file")
			return nil, nil, err
		}
		_, err = io.Copy(outFile, rc)
		// Close the file without defer to close before next iteration of loop
		outFile.Close()
		rc.Close()
		if err != nil {
			fmt.Println("error : close file")
			return nil, nil, err
		}
		filenames = append(filenames, f.Name)
		uuidFileNames = append(uuidFileNames, filename)
	}

	err = os.Remove(path)
	return uuidFileNames, filenames, nil
}

// ReadCsv ...
func ReadCsv(c *gin.Context, csvFilename string) ([][]string, error) {
	path := cfg.Config.Path.Doc
	csvFile, err := c.FormFile(csvFilename)
	if err != nil {
		return nil, fmt.Errorf("Cannot get File")
	}

	if strings.Contains(csvFile.Filename, ".csv") == false {
		return nil, fmt.Errorf("File Must Be CSV")
	}

	filename := GeneratedUUID(filepath.Base(csvFile.Filename))

	if runtime.GOOS == "windows" {
		drive := "D:"
		replacepath := strings.ReplaceAll(path, "/", "\\")
		path = drive + replacepath
	}

	pathCsv := path + filename

	// Upload csv
	if err := c.SaveUploadedFile(csvFile, pathCsv); err != nil {
		return nil, fmt.Errorf("Upload CSV error")
	}
	file, err := os.Open(pathCsv)

	if err != nil {
		return nil, fmt.Errorf("Could Not Open CSV File")
	}

	return csv.NewReader(file).ReadAll()
}

// FindStringInSlice ...
func FindStringInSlice(slice []string, val string) bool {
	for _, item := range slice {
		if item == val {
			return true
		}
	}
	return false
}

// ScheduleCompare ...
func ScheduleCompare(startDate string, endDate string, dates []string, separator string) int {

	var (
		count                     int
		startDateTimestamp        int64
		endDateTimestamp          int64
		startDateCompareTimestamp int64
		endDateCompareTimestamp   int64
	)
	startDateTimestamp = DateStringToMilliSeconds(startDate)
	endDateTimestamp = DateStringToMilliSeconds(endDate)

	for _, value := range dates {
		split := strings.Split(value, separator)
		startDateCompareTimestamp = DateStringToMilliSeconds(split[0])
		endDateCompareTimestamp = DateStringToMilliSeconds(split[1])
		countConvert, _ := strconv.Atoi(split[2])
		if startDateCompareTimestamp == startDateTimestamp && endDateCompareTimestamp == endDateTimestamp {
			continue
		} else {
			if (startDateTimestamp < startDateCompareTimestamp && endDateTimestamp > endDateCompareTimestamp) ||
				(startDateTimestamp >= startDateCompareTimestamp && endDateTimestamp <= endDateCompareTimestamp) ||
				(startDateTimestamp <= startDateCompareTimestamp && endDateTimestamp >= startDateCompareTimestamp &&
					endDateTimestamp <= endDateCompareTimestamp) || (startDateTimestamp >= startDateCompareTimestamp &&
				endDateTimestamp >= endDateCompareTimestamp && startDateTimestamp <= endDateCompareTimestamp) {
				if startDateTimestamp >= startDateCompareTimestamp && endDateTimestamp <= endDateCompareTimestamp {
					count += countConvert
				} else if startDateTimestamp <= startDateCompareTimestamp && endDateTimestamp > startDateCompareTimestamp &&
					endDateTimestamp <= endDateCompareTimestamp {
					count += countConvert
				} else if startDateTimestamp <= startDateCompareTimestamp && startDateTimestamp < endDateCompareTimestamp &&
					endDateTimestamp >= endDateCompareTimestamp {
					count += countConvert
				} else if startDateTimestamp >= startDateCompareTimestamp && endDateTimestamp >= endDateCompareTimestamp &&
					startDateTimestamp < endDateCompareTimestamp {
					count += countConvert
				}
			}
		}

	}
	return count
}

// VisitItemHistory method for formatting visit sales in order history
// format -> dd MMM yyyy. HH:mm - HH:mm
func VisitItemHistory(start, end time.Time) string {
	formatDate := "02 Jan 2006"
	formatTime := "15:04"
	date := start.Format(formatDate)
	time1 := start.Format(formatTime)
	time2 := end.Format(formatTime)
	return fmt.Sprintf("%s : %s - %s", date, time1, time2)
}

// IsExistSlicesInt64 method is item exist in slices int64
func IsExistSlicesInt64(slice *[]int64, item int64) bool {
	if len(*slice) == 0 {
		return false
	}
	for _, a := range *slice {
		if a == item {
			return true
		}
	}
	return false
}
