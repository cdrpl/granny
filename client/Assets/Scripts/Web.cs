using Idlemon.Data;
using System.Collections.Generic;
using System.Net;
using System.Net.Http;
using System.Threading.Tasks;
using UnityEngine;

namespace Idlemon
{
    /// <summary>
    /// Idlemon specific HTTP client.
    /// </summary>
    public class Web
    {
        public static readonly HttpClient Client = new HttpClient();

        /// <summary>
        /// Makes an HTTP request to the Idlemon sign in route.
        /// </summary>
        public static async Task<ApiResponse> SignIn(string email, string pass)
        {
            string url = ApiUrl("/sign-in");

            // Create the form
            var form = new FormUrlEncodedContent(new[] {
                new KeyValuePair<string, string>("email", email),
                new KeyValuePair<string, string>("pass", pass)
            });

            HttpResponseMessage response = await Client.PostAsync(url, form); // Send the request
            string body = await response.Content.ReadAsStringAsync();

            LogResponse(url, body, response.StatusCode);

            return new ApiResponse(body, response.StatusCode);
        }

        public static async Task<ApiResponse> SignUp(string name, string email, string pass)
        {
            string url = ApiUrl("/sign-up");

            // Create the form
            var form = new FormUrlEncodedContent(new[] {
                new KeyValuePair<string, string>("name", name),
                new KeyValuePair<string, string>("email", email),
                new KeyValuePair<string, string>("pass", pass)
            });

            HttpResponseMessage response = await Client.PostAsync(url, form); // Send the request
            string body = await response.Content.ReadAsStringAsync();

            LogResponse(url, body, response.StatusCode);

            return new ApiResponse(body, response.StatusCode);
        }

        public static async Task<string> GetRoom()
        {
            string url = ServerUrl("/room");

            HttpResponseMessage response = await Client.GetAsync(url);
            response.EnsureSuccessStatusCode();
            string body = await response.Content.ReadAsStringAsync();

            LogResponse(url, body, response.StatusCode);

            return body;
        }

        /// <summary>
        /// Returns the full URL including the path for the public API.
        /// </summary>
        public static string ApiUrl(string path)
        {
            return Const.WEB_PROTOCOL + "://" + Const.API_ADDR + path;
        }

        static string ServerUrl(string path)
        {
            return Const.WEB_PROTOCOL + "://" + Const.SERVER_ADDR + path;
        }

        static void LogResponse(string path, string body, HttpStatusCode statusCode)
        {
            Debug.Log(path + " " + statusCode.ToString() + " " + body);
        }
    }
}
