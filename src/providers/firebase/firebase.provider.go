package firebase

import (
	"context"
	firebase "firebase.google.com/go/v4"
	"go.uber.org/fx"
	"google.golang.org/api/option"
	"log"
)

type FirebaseProvider struct {
	GetApp func() *firebase.App
}

func NewFirebaseProvider() *FirebaseProvider {
	opt := option.WithCredentialsFile("firebase-adminsdk.json")
	config := &firebase.Config{ProjectID: "myapp-39549"}
	app, err := firebase.NewApp(context.Background(), config, opt)
	if err != nil {
		log.Fatalf("error initializing app: %v\n", err)
	}
	return &FirebaseProvider{
		GetApp: func() *firebase.App {
			return app
		},
	}
}

var Module = fx.Module("FirebaseProvider",
	fx.Provide(NewFirebaseProvider),
	fx.Invoke(func(firebaseProvider *FirebaseProvider) {
		log.Printf("FirebaseProvider: %v\n", firebaseProvider)
	}))
