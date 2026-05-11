package main

import (
    "log"
    "net/http"
)

func main() {
    // Servir le build React (frontend/dist) à la racine
    http.Handle("/", http.FileServer(http.Dir("./frontend/dist")))

    // Servir les fichiers WASM sous /wasm/
    http.Handle("/wasm/", http.StripPrefix("/wasm/", http.FileServer(http.Dir("./static"))))

    log.Println("Serveur sur http://localhost:8080")
    log.Fatal(http.ListenAndServe(":8080", nil))
}