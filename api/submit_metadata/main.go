package submit_metadata

import (
	"github.com/aws/aws-lambda-go/events"
	"log"
	"github.com/codetaming/indy-ingest/api/validator"
	"encoding/json"
	"github.com/aws/aws-lambda-go/lambda"
)

func Handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	log.Printf("Processing Lambda request %s\n", request.RequestContext.RequestID)

	schemaUrl := request.Headers["describedBy"]
	bodyJson := request.Body
	result := validator.Validate(schemaUrl, bodyJson)
	body, _ := json.Marshal(result)

	if result.Valid {
		return events.APIGatewayProxyResponse{
			Headers:    map[string]string{"Content-Type": "application/json"},
			Body:       string(body),
			StatusCode: 200,
		}, nil
	} else
	{
		return events.APIGatewayProxyResponse{
			Headers:    map[string]string{"Content-Type": "application/json"},
			Body:       string(body),
			StatusCode: 400,
		}, nil
	}

}

func main() {
	lambda.Start(Handler)
}
