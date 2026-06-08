package logger

import (
	"log"
	
)

var debugMode = false

// EnableDebug active les logs de debug
func EnableDebug() {
	debugMode = true
	log.Println("🔍 Mode DEBUG activé")
}

// DisableDebug désactive les logs de debug
func DisableDebug() {
	debugMode = false
}

// Debugf affiche un log de debug si le mode est activé
func Debugf(format string, v ...interface{}) {
	if debugMode {
		log.Printf("[DEBUG] "+format, v...)
	}
}

// Infof affiche toujours un log d'info
func Infof(format string, v ...interface{}) {
	log.Printf("[INFO] "+format, v...)
}

// Errorf affiche toujours un log d'erreur
func Errorf(format string, v ...interface{}) {
	log.Printf("[ERROR] "+format, v...)
}