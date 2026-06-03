package runner

import (
	"fmt"
	"net/http"
	"os"

	xtremecore "github.com/globalxtreme/go-core/v2"
	xtremepkg "github.com/globalxtreme/go-core/v2/pkg"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"github.com/spf13/cobra"

	"service/internal/app/api"
	"service/internal/pkg/config"
)

var rootCmd = &cobra.Command{
	Use:  "root",
	Long: "Running service api",
	Run: func(cmd *cobra.Command, args []string) {
		xtremepkg.InitDevMode()
		xtremepkg.InitHost()

		config.InitTZ()
		config.InitCors()
		config.InitMail()
		config.InitRPC()
		config.InitValidation()

		DBClose := config.InitDB()
		defer DBClose()

		// TODO: Aktifkan saat up ke operational
		//xtremedb.Migrate(config.PgSQL, database.Migrations())

		// rabbitMQClose := config.InitRabbitMQ()
		// defer rabbitMQClose()

		logCleanup := xtremepkg.InitLogRPC()
		defer logCleanup()

		newCors := cors.New(config.CorsOptions)

		newRoute := mux.NewRouter()
		xtremecore.RegisterRouter(newRoute, api.Register)

		fmt.Println(fmt.Sprintf("Server started on %s", xtremepkg.HostFull))

		err := http.ListenAndServe(xtremepkg.Host, newCors.Handler(newRoute))
		if err != nil {
			panic(err)
		}
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().BoolVar(&xtremepkg.DevMode, "dev", false, "Set for development mode")
}
