package authtransport

import (
	"net/http"

	"github.com/Andryshazzz/go_pet/internal/core/domain"
	logger "github.com/Andryshazzz/go_pet/internal/core/logger"
	httprequest "github.com/Andryshazzz/go_pet/internal/core/transport/http/request"
	httpresponse "github.com/Andryshazzz/go_pet/internal/core/transport/http/response"
)

// CreateUserRequest represents the JSON body for creating a new user.
type CreateUserRequest struct {
	FullName    string  `json:"full_name" validate:"required,min=2,max=100"                     example:"Ivan Ivanov"`
	PhoneNumber *string `json:"phone_number" validate:"omitempty,min=10,max=15,startswith=+"    example:"+79991118899"`
}

// CreateUserResponse represents the JSON response after creating a user.
type CreateUserResponse struct {
	ID          string  `json:"id"`
	FullName    string  `json:"full_name"`
	PhoneNumber *string `json:"phone_number"`
}

// CreateUser	godoc
// @Summary  	Создать пользователя
// @Description Создать нового пользователя в системе
// @Tags 		Auth
// @Accept  	json
// @Produce 	json
// @Param   	request body CreateUserRequest true "CreateUser тело запроса"
// @Success 	201 {object} CreateUserResponse "Успешно созданный пользователь"
// @Failure 	400 {object} httpresponse.ErrorsResponse "Bad request"
// @Failure 	500 {object} httpresponse.ErrorsResponse "Internal server error"
// @Router		/users [post]
func (h *UsersHTTPHandler) CreateUser(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := logger.FromContext(ctx)
	responseHandler := httpresponse.NewHTTPResponseHandler(log, rw)

	var request CreateUserRequest

	if err := httprequest.DecodeAndValidateRequest(r, &request); err != nil {
		responseHandler.ErrorResponse(err, "failed to decode and validate HTTP request")
		return
	}

	userDomain := domainFromDTO(request)

	userDomain, err := h.usersService.CreateUser(ctx, userDomain)

	if err != nil {
		responseHandler.ErrorResponse(err, "failed to create user")

		return
	}

	response := dtoFromDomain(userDomain)

	responseHandler.JSONResponse(response, http.StatusCreated)

}

// domainFromDTO maps a CreateUserRequest (DTO) to a domain User entity.
func domainFromDTO(dto CreateUserRequest) domain.User {
	return domain.NewUninitializedUser(dto.FullName, dto.PhoneNumber)
}

// dtoFromDomain maps a domain User entity to a CreateUserResponse (DTO).
func dtoFromDomain(user domain.User) CreateUserResponse {
	return CreateUserResponse{
		ID:          user.ID,
		FullName:    user.FullName,
		PhoneNumber: user.PhoneNumber,
	}
}
