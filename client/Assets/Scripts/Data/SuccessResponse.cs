using Newtonsoft.Json.Linq;

namespace Idlemon.Data
{
    public class SuccessResponse
    {
        public JToken Data { get; private set; }

        public SuccessResponse(JToken data)
        {
            Data = data;
        }
    }
}
