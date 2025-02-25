package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/zmb3/spotify"
)

var (
	auth  = spotify.NewAuthenticator("http://localhost:8080/callback", spotify.ScopePlaylistReadPrivate)
	ch    = make(chan *spotify.Client)
	state = "random-string-for-security"
)

func main() {
	clientID := "d94076b3d4174dfca04516c874590d20"
	clientSecret := "8cb5a32df82541d59438409bbcb79bb5"

	auth.SetAuthInfo(clientID, clientSecret)

	http.HandleFunc("/callback", completeAuth)
	go http.ListenAndServe(":8080", nil)

	url := auth.AuthURL(state)
	fmt.Printf("Veuillez ouvrir le lien suivant dans votre navigateur pour authentification :\n%s\n", url)

	client := <-ch

	user, err := client.CurrentUser()
	if err != nil {
		log.Fatalf("Impossible de récupérer les informations de l'utilisateur : %v", err)
	}
	fmt.Printf("Connecté en tant que : %s (%s)\n", user.DisplayName, user.ID)

	playlists, err := client.CurrentUsersPlaylists()
	if err != nil {
		log.Fatalf("Impossible de récupérer les playlists : %v", err)
	}

	fmt.Println("Playlists :")
	for _, playlist := range playlists.Playlists {
		fmt.Printf("- %s\n", playlist.Name)
	}
}

func completeAuth(w http.ResponseWriter, r *http.Request) {
	tok, err := auth.Token(state, r)
	if err != nil {
		http.Error(w, "Impossible d'obtenir le token", http.StatusForbidden)
		log.Fatalf("Impossible d'obtenir le token : %v", err)
		return
	}
	if st := r.FormValue("state"); st != state {
		http.NotFound(w, r)
		log.Fatalf("État invalide")
		return
	}
	client := auth.NewClient(tok)
	fmt.Fprintln(w, "Authentification réussie ! Vous pouvez revenir à l'application.")
	ch <- &client
}
