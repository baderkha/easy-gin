# easy-gin

<p align="center"> <img src="./logo/logo.png" width="200px" height="200px"/> </p>

A simpler way to make [Gin  REST HTTP handlers](https://github.com/gin-gonic/gin) following the [DTO pattern](https://www.okta.com/identity-101/dto/) for the requests. 
The goal is to have handlers that are easily unit testable for restful api projects . 

Example : 
- Old Fashion way : bad seperation of concern not easily testable because of the context object

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
 - Easy gin way : seperation of concern and easy input output with your own objects , making it easy to test !
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

### supported 
