package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/Shulammite-Aso/bazel-demo-app/bazel"
	"github.com/Shulammite-Aso/bazel-demo-app/handlers"
	"github.com/antchfx/xmlquery"
	"github.com/bgentry/go-netrc/netrc"
	"github.com/bwmarrin/snowflake"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/fatih/color"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	cache "github.com/patrickmn/go-cache"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
	"gopkg.in/yaml.v3"
)

// Config struct for demonstration with validator tags
type Config struct {
	AppName string `yaml:"app_name" validate:"required"`
	Port    int    `yaml:"port" validate:"required,min=1000,max=65535"`
	Debug   bool   `yaml:"debug"`
}

// Simple protobuf message demonstration
type ProtoMessage struct {
	Timestamp *timestamppb.Timestamp
	Message   string
}

func doNotInvoke() (string, error) {
	return bazel.Runfile("/tmp/does/not/exist")
}

// demonstrateNewDependencies shows usage of all new dependencies
func demonstrateNewDependencies() {
	// 1. godotenv - Load environment variables
	_ = godotenv.Load() // Load from .env file if it exists
	color.Green("✓ godotenv: Environment loaded")

	// 2. viper - Configuration management
	viper.SetDefault("app_name", "bazel-demo-app")
	viper.SetDefault("port", 5000)
	viper.SetDefault("debug", true)
	color.Green("✓ viper: Configuration set with defaults")

	// 3. yaml.v3 - YAML parsing
	config := Config{
		AppName: viper.GetString("app_name"),
		Port:    viper.GetInt("port"),
		Debug:   viper.GetBool("debug"),
	}
	yamlData, _ := yaml.Marshal(&config)
	color.Green("✓ yaml.v3: Config marshaled to YAML: %s", string(yamlData))

	// 4. validator - Struct validation
	validate := validator.New()
	if err := validate.Struct(config); err == nil {
		color.Green("✓ validator: Config validation passed")
	}

	// 5. go-cache - In-memory caching
	c := cache.New(5*time.Minute, 10*time.Minute)
	c.Set("demo-key", "demo-value", cache.DefaultExpiration)
	if val, found := c.Get("demo-key"); found {
		color.Green("✓ go-cache: Retrieved value from cache: %s", val)
	}

	// 6. jwt-go - JWT token generation
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user": "demo-user",
		"exp":  time.Now().Add(time.Hour * 24).Unix(),
	})
	tokenString, _ := token.SignedString([]byte("secret-key"))
	color.Green("✓ jwt-go: Generated JWT token (truncated): %s...", tokenString[:20])

	// 7. testify/assert - Assertions (typically for tests, but demonstrating here)
	testValue := true
	assert.True(nil, testValue, "This should be true") // nil context for demo
	color.Green("✓ testify/assert: Assertion passed")

	// 8. protobuf - Protocol buffers
	protoMsg := &ProtoMessage{
		Timestamp: timestamppb.Now(),
		Message:   "Hello from protobuf",
	}
	// Marshal to demonstrate protobuf usage
	_ = proto.Size(protoMsg.Timestamp)
	color.Green("✓ protobuf: Created protobuf timestamp: %v", protoMsg.Timestamp.AsTime())

	// 10. fatih/color - Already used above for colored output
	color.Yellow("✓ color: All dependencies demonstrated successfully!")

	// 11. cobra - CLI framework (command structure)
	color.Cyan("✓ cobra: CLI framework initialized (see rootCmd)")
}

var rootCmd = &cobra.Command{
	Use:   "bazel-demo-app",
	Short: "A demo Bazel Go application",
	Long:  "A demonstration application showing Bazel build with multiple Go dependencies",
	Run: func(cmd *cobra.Command, args []string) {
		runServer()
	},
}

func runServer() {
	fmt.Println("Hello world")

	// Demonstrate all new dependencies
	color.Cyan("\n=== Demonstrating New Dependencies ===")
	demonstrateNewDependencies()
	color.Cyan("=====================================\n")

	// Existing functionality
	wadl, err := xmlquery.LoadURL("https://httpbin.org/get")
	if err != nil {
		panic(err)
	}

	netrc := netrc.Machine{
		Login:    "test",
		Password: "test",
		Account:  "test",
	}

	fmt.Println(netrc)

	logrus.Info("Hello world")

	uuid, err := uuid.NewUUID()
	if err != nil {
		panic(err)
	}

	fmt.Println(uuid)

	sf, err := snowflake.NewNode(1)
	if err != nil {
		panic(err)
	}

	fmt.Println(sf.Generate())

	attr := xmlquery.FindOne(wadl, "//application/@xmlns")
	fmt.Println(attr.InnerText())

	router := mux.NewRouter()

	router.HandleFunc("/greet", handlers.Greet).Methods("GET")
	router.HandleFunc("/greet-many", handlers.GreetMany).Methods("GET")

	address := ":5000"

	log.Printf("server started at port %v\n", address)

	err = http.ListenAndServe(address, router)

	if errors.Is(err, http.ErrServerClosed) {
		log.Printf("server closed\n")
	} else if err != nil {
		log.Printf("error starting server: %s\n", err)
		os.Exit(1)
	}
}

func main() {
	// Use cobra for CLI command handling
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
