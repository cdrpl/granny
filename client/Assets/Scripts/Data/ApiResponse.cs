using Newtonsoft.Json;
using Newtonsoft.Json.Linq;
using System.Net;

namespace Idlemon.Data
{
    public class ApiResponse
    {
        public SuccessResponse Success { get; private set; }
        public ErrorResponse Error { get; private set; }
        public HttpStatusCode StatusCode { get; private set; }

        public bool HasData => Success != null;
        public bool HasError => Error != null;

        public ApiResponse(string response, HttpStatusCode statusCode)
        {
            JObject obj = JObject.Parse(response);

            JToken err = obj["error"];
            JToken data = obj["data"];

            if (err != null)
            {
                Error = JsonConvert.DeserializeObject<ErrorResponse>(err.ToString());
            }
            else if (data != null)
            {
                Success = new SuccessResponse(data);
            }

            this.StatusCode = statusCode;
        }
    }
}
