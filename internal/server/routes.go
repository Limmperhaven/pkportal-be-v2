package server

import (
	"github.com/Limmperhaven/pkportal-be-v2/internal/controllers"
	"github.com/Limmperhaven/pkportal-be-v2/internal/controllers/middlewares"
	"github.com/gin-gonic/gin"
	"github.com/rs/cors"
	"net/http"
)

func initRoutes(router *gin.Engine, c *controllers.ControllerStorage, m *middlewares.MiddlewareStorage) {
	authorized := router.Group("/")
	authorized.Use(m.AuthMiddleware)
	//authorized.Use(m.CheckActivationMiddleware)
	admin := authorized.Group("/")
	admin.Use(m.CheckAdminRoleMiddleware)

	auth := router.Group("/auth")
	{
		auth.POST("/sign-up", c.SignUp)
		auth.POST("/sign-in", c.SignIn)
		auth.POST("/recover/:email", c.RecoverPassword)
		auth.GET("/activate/:token", c.ActivateAccount)
		auth.POST("/logout", c.Logout)
		auth.POST("/confirmRecover/:token", c.ConfirmRecover)
	}

	authTd := authorized.Group("/td")
	adminTd := admin.Group("/td")
	{
		adminTd.POST("/create", c.CreateTestDate)
		authTd.GET("/byId/:id", c.GetTestDate)
		adminTd.PUT("/byId/:id", c.UpdateTestDate)
		adminTd.POST("/setStatus/:id/:status", c.SetTestDatePubStatus)
		adminTd.GET("/list", c.ListTestDates)
		authTd.GET("/listAvailable", c.ListAvailableTestDates)
		adminTd.POST("/signUpUser/:userId/:tdId", c.SignUpUserToTestDate)
		adminTd.POST("/signUpMe/:tdId", c.SignUpMeToTestDate)
		adminTd.GET("/listCommonLocations", c.ListCommonLocations)
		adminTd.POST("/setAttended/:userId/:tdId", c.SetTestDateAttended)
	}

	authUser := authorized.Group("/user")
	adminUser := admin.Group("/user")
	{
		adminUser.POST("/create", c.CreateUser)
		adminUser.GET("/byId/:id", c.GetUser)
		adminUser.PUT("/byId/:id", c.UpdateUser)
		adminUser.GET("/list", c.ListUsers)
		adminUser.POST("/setStatus/:userId/:statusId", c.SetUserStatus)
		adminUser.GET("/downloadScreenshot/:userId", c.DownloadScreenshot)
		authUser.GET("/me", c.GetMe)
		authUser.GET("/listStatuses", c.ListStatuses)
		authUser.POST("/uploadScreenshot", c.UploadScreenshot)
	}

	authProfile := authorized.Group("/profiles")
	adminProfile := admin.Group("/profiles")
	{
		adminProfile.POST("/create", c.CreateProfile)
		authProfile.GET("/byId/:id", c.CreateProfile)
		adminProfile.PUT("/byId/:id", c.UpdateProfile)
		authProfile.GET("/list", c.ListProfiles)
		adminProfile.POST("/setToUser/:userId", c.SetProfilesToUser)
		authProfile.POST("/setToMe", c.SetProfilesToMe)
	}

	authSubject := authorized.Group("/subjects")
	adminSubject := admin.Group("/subjects")
	{
		adminSubject.POST("/create", c.CreateSubject)
		authSubject.GET("/byId/:id", c.GetSubject)
		adminSubject.PUT("/byId/:id", c.UpdateSubject)
		authSubject.GET("/list", c.ListSubjects)
		authSubject.GET("/listFL", c.ListForeignLanguages)
		adminSubject.POST("/setToUser/:userId", c.SetSubjectToUser)
		authSubject.POST("/setToMe", c.SetSubjectToMe)
	}

	authFL := authorized.Group("/fl")
	adminFL := admin.Group("/fl")
	{
		adminFL.POST("/create", c.CreateForeignLanguage)
		authFL.GET("/byId/:id", c.GetForeignLanguage)
		adminFL.POST("/byId/:id", c.UpdateForeignLanguage)
		authFL.GET("/list", c.ListForeignLanguages)
		adminFL.POST("/setToUser/:userId/:flId", c.SetForeignLanguageToUser)
		authFL.POST("/setToMe/:flId", c.SetForeignLanguageToMe)
	}
}

func initCors(router *gin.Engine) http.Handler {
	c := cors.New(cors.Options{
		AllowOriginFunc:        func(origin string) bool { return true },
		AllowOriginRequestFunc: func(r *http.Request, origin string) bool { return true },
		AllowedMethods: []string{
			http.MethodGet, http.MethodPost, http.MethodPut,
		},
		AllowedHeaders:      []string{"accept", "authorization", "content-type"},
		ExposedHeaders:      []string{"Set-Cookie", "authorization", "Content-Disposition"},
		AllowCredentials:    true,
		AllowPrivateNetwork: true,
		OptionsPassthrough:  false,
		Debug:               true,
	})
	return c.Handler(router)
}
