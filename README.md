# easy-gin

<p align="center"> <img src="./logo/logo.png" width="200px" height="200px"/> </p>

A No magic simple way to implement [Gin  REST HTTP handlers](https://github.com/gin-gonic/gin) following the [DTO pattern](https://www.okta.com/identity-101/dto/) . 
The goal is to have handlers that are unit testable and free of boilerplate code. 

**NOTE**
This package is meant for REST APIs only using the Gin framework.

Example : 
- Without easy-gin : 
	- bad seperation of concern 
		- doing validation and business logic
	- not easily unit testable without having to either wrap the context or use an http recorder
	- repetitive code with binds and checking error ..etc
	

  ```go
   type UserInput struct {
      UserID            string `json:"user_id" uri:"user_id"`
    }

    func main(){
      en := gin.Default()
      en.POST("/:user_id", func(ctx *gin.Context) {
          var u UserInput
          // bind
          err := ctx.BindQuery(&u)
          if err != nil {
            ctx.JSON(http.StatusBadRequest, err.Error())
            return
          }
          // some extra validation logic
          // i know you can add rules in the struct , but this can be replaced with a validation from db ...etc
          if u.UserID == "" {
            ctx.JSON(http.StatusBadRequest, errors.New("user id is missing"))
            return
          }

          // do somethign with data

          ctx.JSON(http.StatusOK, fmt.Sprintf("user with id %s has been processed", u.UserID))

	    })
    }
  ```
 - easy-gin way : 
 	- unit testable
 		- your test is input output 
 		- no need to mock or use external recorders ..etc 
 	- seperation of concern 
 		- validation is handled by the UserInput DTO
 		- Business Logic is handled by the handler  


  ``` go
    var _ easygin.IRequest = &UserInput{} // this struct implements IRequest 
    type UserInput struct {
      UserID            string `json:"user_id" uri:"user_id"` // still use the bind methods from gin !
    }
    // add custom validation logic not restricted by struct tags
    func (u UserInput) Validate() error {
      if u.UserID == "" {
        return errors.New("user id not set")
      }
      return nil
    }
    // you can use this to wrap your error if validation failed from gin or your custom validation
    func (u UserInput) ValidationErrorFormat(err error) any {
      return map[string]any{
        "err": err.Error(),
      }
    }
    
    func HandleUsers(u UserInput) *easygin.Response {
      // do something with the input ...
      // focus on your domain logic rather than validation ...etc
      return easygin.
        Res(fmt.Sprintf("user with id %s has been processed",u.UserID)).
        Status(http.StatusOK)
    }
    
    func main() {
      en := gin.Default()
      // by default the second argument is optional 
      // if not provided it will atempt all bind methods (JSON,QUERY,URI) (this will incur a performance hit)
      en.POST("/:user_id", easygin.To(HandleUsers,easygin.BindURI)) 
      en.Run(":80")
    }
  ```
 
## Installation 
*Note* you need atleast golang 1.8 and above to install this utility as under the hood it uses generics
```
go get -u github.com/baderkha/easy-gin/v1/easygin
```

## Documentation

### Quick Setup

- Step 1 : Create a DTO object that implements the easygin.IRequest interface 

	```go 
	type UpdateUserRequest struct {
		ID string `uri:"user_id"` // regular gin binding from instructions
		Name string `json:"name"` // binds from json body
		UserType string `form:"user_type"` // binds from query parameter
	}
	
	func (u UpdateUserRequest) Validate() error {
	   // your custom validation here , consider using a struct validator like [go-validator](https://github.com/go-playground/validator)
	   // also you can just have your validation done via tags if it's simple stuff and just return nil here
	   return nil
	}
	
	func (u UpdateUserRequest) ValidationErrorFormat(err error) any {
		// if you want your response to be the error string return
		return err.Error()
		// if you want your response to be a wrapped with an object (map option)
		return map[string]any{
			"err":err.Error()
		}
		// if you want your response to be a wrapped with an object (struct option)
		return struct {
			Error   string `json:"err"`
			Message string `json:"server_message"`
		}{
			Error:   err.Error(),
			Message: "failed validation",
		}
	}
	```
- Step 2 : Create your easygin Handler
	``` go
	// argument must not be a pointer !
	func HandleUserUpdate(u UpdateUserRequest) *easygin.Response {
		// process the data 
		// ....
		
		// once ready to respond to client
		res := easygin.Res(map[string]any{"wow":"ok"})
		
		return res // this will default with a 200 response code 
		
		return res.Status(201) // you can override it yourself , so you can use this to handle errors 
	}
	```
- Step 3 : Add it to your routes
	``` go
	func main() {
		en := gin.Default()
		
		// option a default binding
		// although this looks cleaner this will have a performance hit if you do not need to bind from everything else
		en.PATCH("/:user_id",easygin.To(HandleUserUpdate)) 
		
		// option b recommended
		// only bind from ...
		// preferable  ,always define where you're binding from
		en.PATCH("/:user_id",easygin.To(HandleUserUpdate,easygin.BindURI,easygin.BindJSON,easygin.BindQuery))
	}
		
	```

That's it , this should now work and bind from all the different part of the http request
