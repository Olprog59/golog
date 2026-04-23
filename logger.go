package golog

import (
	"encoding/json"
	"fmt"
	"log"
	"log/syslog"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"
)

var (
	errColor     = makeColor("\033[41;97m")  // fond rouge, texte blanc
	warnColor    = makeColor("\033[43;30m")  // fond jaune, texte noir
	debugColor   = makeColor("\033[44;97m")  // fond bleu, texte blanc
	infoColor    = makeColor("\033[45;97m")  // fond magenta, texte blanc
	noticeColor  = makeColor("\033[46;30m")  // fond cyan, texte noir
	successColor = makeColor("\033[42;30m")  // fond vert, texte noir
)

// mu protects all mutable package-level state.
var mu sync.RWMutex

var (
	syslogWriter    *syslog.Writer
	separator       = " | "
	timeFormat      string
	language        = "en"
	includeFileName bool
)

// timePrecision représente la précision des secondes dans le format horaire.
type timePrecision string

const (
	MILLI timePrecision = ".000"
	MICRO timePrecision = ".000000"
	NANO  timePrecision = ".000000000"
)

type LogMessage struct {
	Timestamp string `json:"timestamp"`
	Level     string `json:"level"`
	FileLine  string `json:"fileLine,omitempty"`
	Message   string `json:"message"`
	SessionID string `json:"sessionID,omitempty"`
	UserID    string `json:"userID,omitempty"`
}

// colorFunc est le type des fonctions de formatage avec couleur.
type colorFunc func(message, level, fileLine, tf, sep string, sessionID, userID *string, args ...any) string

func SetLanguage(lang string) {
	mu.Lock()
	defer mu.Unlock()
	switch lang {
	case "en":
		timeFormat = "01-02-2006 3:04:05 PM"
	case "fr":
		timeFormat = "02-01-2006 15:04:05"
	case "de":
		timeFormat = "02.01.2006 15:04:05"
	case "zh", "ja":
		timeFormat = "2006年01月02日 15:04:05"
	case "ru":
		timeFormat = "02.01.2006 15:04:05"
	case "es", "it":
		timeFormat = "02-01-2006 15:04:05"
	case "ar":
		timeFormat = "02/01/2006 15:04:05"
	default:
		timeFormat = "2006-01-02 15:04:05"
	}
	language = lang
}

// InitSyslog initialise le writer syslog distant.
func InitSyslog(network, raddr, nameApp string) error {
	w, err := syslog.Dial(network, raddr, syslog.LOG_INFO|syslog.LOG_LOCAL0, nameApp)
	if err != nil {
		return err
	}
	mu.Lock()
	syslogWriter = w
	mu.Unlock()
	return nil
}

// SetTimePrecision ajuste la précision des secondes dans le format horaire.
func SetTimePrecision(precision timePrecision) {
	mu.Lock()
	defer mu.Unlock()
	if language == "en" {
		timeFormat = "01-02-2006 3:04:05" + string(precision) + " PM"
	} else {
		timeFormat = timeFormat + string(precision)
	}
}

// SetCustomTimeFormat applique un format horaire personnalisé.
func SetCustomTimeFormat(customFormat string) {
	mu.Lock()
	timeFormat = customFormat
	mu.Unlock()
}

// SetSeparator définit le séparateur entre les champs du message de log.
func SetSeparator(sep string) {
	mu.Lock()
	separator = sep
	mu.Unlock()
}

// EnableFileNameLogging active l'affichage du fichier et de la ligne dans les logs.
func EnableFileNameLogging() {
	mu.Lock()
	includeFileName = true
	mu.Unlock()
}

// getFileAndLine retourne "fichier:ligne" de l'appelant. Profondeur 3 : cette
// fonction → format → Err/Info/… → code utilisateur.
func getFileAndLine() string {
	_, file, line, ok := runtime.Caller(3)
	if !ok {
		return "unknown:0"
	}
	return file[strings.LastIndex(file, "/")+1:] + ":" + strconv.Itoa(line)
}

// makeColor retourne une colorFunc dont le badge de niveau (uniquement) est coloré
// avec ansiOpen. Le reste du message s'affiche dans la couleur par défaut du terminal,
// garantissant la lisibilité quel que soit le fond (sombre ou clair).
func makeColor(ansiOpen string) colorFunc {
	return func(message, level, fileLine, tf, sep string, sessionID, userID *string, args ...any) string {
		formattedMessage := fmt.Sprintf(message, args...)
		timestamp := time.Now().Format(tf)

		var b strings.Builder
		b.WriteString(timestamp)
		b.WriteString(sep)
		b.WriteString(ansiOpen)
		fmt.Fprintf(&b, " %-7s", level)
		b.WriteString("\033[0m")

		if fileLine != "" {
			b.WriteString(sep)
			fmt.Fprintf(&b, "%-15s", fileLine)
		}

		var sID, uID *string
		if sessionID != nil && *sessionID != "" {
			sID = sessionID
		}
		if userID != nil && *userID != "" {
			uID = userID
		}

		switch {
		case sID != nil && uID != nil:
			b.WriteString(sep)
			b.WriteString("sessionID: ")
			b.WriteString(*sID)
			b.WriteString(sep)
			b.WriteString("userID: ")
			b.WriteString(*uID)
		case sID != nil:
			b.WriteString(sep)
			b.WriteString("sessionID: ")
			b.WriteString(*sID)
		case uID != nil:
			b.WriteString(sep)
			b.WriteString("userID: ")
			b.WriteString(*uID)
		}

		b.WriteString(sep)
		b.WriteString(formattedMessage)

		return b.String()
	}
}

func enrichMessage(level, message, fileLine, tf, sessionID, userID string) string {
	logMsg := LogMessage{
		Timestamp: time.Now().Format(tf),
		Level:     level,
		FileLine:  fileLine,
		Message:   message,
		SessionID: sessionID,
		UserID:    userID,
	}

	jsonMessage, err := json.Marshal(logMsg)
	if err != nil {
		fmt.Println("Erreur lors de l'encodage JSON du message:", err)
		return message
	}

	return string(jsonMessage)
}

// format est le point d'entrée commun à toutes les fonctions de log.
// Il lit la configuration globale une seule fois sous verrou, calcule fileLine
// une seule fois si nécessaire, puis délègue à la colorFunc et (conditionnellement)
// au syslog.
func format(fn colorFunc, message, level string, sessionID, userID *string, params ...any) {
	mu.RLock()
	sw := syslogWriter
	inc := includeFileName
	tf := timeFormat
	sep := separator
	mu.RUnlock()

	fileLine := ""
	if inc || sw != nil {
		fileLine = getFileAndLine()
	}

	// Pour l'affichage terminal : fileLine uniquement si EnableFileNameLogging.
	displayLine := ""
	if inc {
		displayLine = fileLine
	}

	fmt.Println(fn(message, level, displayLine, tf, sep, sessionID, userID, params...))

	if sw != nil {
		var sid, uid string
		if sessionID != nil {
			sid = *sessionID
		}
		if userID != nil {
			uid = *userID
		}
		logToSyslog(sw, syslogCheckLevel(level), enrichMessage(level, message, fileLine, tf, sid, uid))
	}
}

func Err(message string, params ...any) {
	format(errColor, message, "Err", nil, nil, params...)
}

func Warn(message string, params ...any) {
	format(warnColor, message, "Warn", nil, nil, params...)
}

func Debug(message string, params ...any) {
	format(debugColor, message, "Debug", nil, nil, params...)
}

func Info(message string, params ...any) {
	format(infoColor, message, "Info", nil, nil, params...)
}

func Notice(message string, params ...any) {
	format(noticeColor, message, "Notice", nil, nil, params...)
}

func Success(message string, params ...any) {
	format(successColor, message, "Success", nil, nil, params...)
}

func ErrorWithID(message, sessionID, userID string, params ...any) {
	format(errColor, message, "Err", &sessionID, &userID, params...)
}

func WarnWithID(message, sessionID, userID string, params ...any) {
	format(warnColor, message, "Warn", &sessionID, &userID, params...)
}

func DebugWithID(message, sessionID, userID string, params ...any) {
	format(debugColor, message, "Debug", &sessionID, &userID, params...)
}

func InfoWithID(message, sessionID, userID string, params ...any) {
	format(infoColor, message, "Info", &sessionID, &userID, params...)
}

func NoticeWithID(message, sessionID, userID string, params ...any) {
	format(noticeColor, message, "Notice", &sessionID, &userID, params...)
}

func SuccessWithID(message, sessionID, userID string, params ...any) {
	format(successColor, message, "Success", &sessionID, &userID, params...)
}

// logToSyslog prend le writer en paramètre pour éviter une lecture globale supplémentaire.
func logToSyslog(w *syslog.Writer, level syslog.Priority, message string) {
	var err error
	switch level {
	case syslog.LOG_ERR:
		err = w.Err(message)
	case syslog.LOG_WARNING:
		err = w.Warning(message)
	case syslog.LOG_INFO:
		err = w.Info(message)
	case syslog.LOG_NOTICE:
		err = w.Notice(message)
	case syslog.LOG_DEBUG:
		err = w.Debug(message)
	default:
		err = w.Info(message)
	}
	if err != nil {
		log.Println(err)
	}
}

func syslogCheckLevel(level string) syslog.Priority {
	switch level {
	case "Err":
		return syslog.LOG_ERR
	case "Warn":
		return syslog.LOG_WARNING
	case "Debug":
		return syslog.LOG_DEBUG
	case "Info":
		return syslog.LOG_INFO
	case "Notice":
		return syslog.LOG_NOTICE
	default:
		return syslog.LOG_INFO
	}
}
