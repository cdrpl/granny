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
        /// Makes an HTTP request to the Idlemon health check route.
        /// </summary>
        public static async Task<string> HealthCheck()
        {
            string url = PublicUrl("");

            HttpResponseMessage response = await Client.GetAsync(url);
            response.EnsureSuccessStatusCode();
            string body = await response.Content.ReadAsStringAsync();

            LogResponse(url, body, response.StatusCode);

            return await response.Content.ReadAsStringAsync();
        }

        /// <summary>
        /// Makes an HTTP request to the Idlemon sign in route.
        /// </summary>
        public static async Task<ApiResponse> SignIn(string email, string pass)
        {
            string url = PublicUrl("/sign-in");

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
            string url = PublicUrl("/sign-up");

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

        public static async Task<ApiResponse> GetRooms()
        {
            HttpResponseMessage response = await Client.GetAsync(LobbyUrl("/rooms"));
            response.EnsureSuccessStatusCode();

            string body = await response.Content.ReadAsStringAsync();

            LogResponse(LobbyUrl("/rooms"), body, response.StatusCode);

            return new ApiResponse(body, response.StatusCode);
        }

        public static async Task<ApiResponse> CreateRoom(int userId, string name)
        {
            string url = LobbyUrl("/rooms");

            var form = new FormUrlEncodedContent(new[] {
                new KeyValuePair<string, string>("userId", userId.ToString()),
                new KeyValuePair<string, string>("roomName", name),
            });

            HttpResponseMessage response = await Client.PostAsync(url, form);
            string body = await response.Content.ReadAsStringAsync();

            LogResponse(url, body, response.StatusCode);

            return new ApiResponse(body, response.StatusCode);
        }

        public static async Task<ApiResponse> JoinRoom(string roomId)
        {
            string url = LobbyUrl("/rooms/join");

            var form = new FormUrlEncodedContent(new[] {
                new KeyValuePair<string, string>("roomId", roomId),
            });

            HttpResponseMessage response = await Client.PostAsync(url, form);
            string body = await response.Content.ReadAsStringAsync();

            LogResponse(url, body, response.StatusCode);

            return new ApiResponse(body, response.StatusCode);
        }

        /// <summary>
        /// Returns the full URL including the path for the public API.
        /// </summary>
        public static string PublicUrl(string path)
        {
            return Const.WEB_PROTOCOL + "://" + Const.PUBLIC_API + path;
        }

        /// <summary>
        /// Returns the full URL including the path for the private API.
        /// </summary>
        public static string LobbyUrl(string path)
        {
            return Const.WEB_PROTOCOL + "://" + Const.LOBBY_ADDR + path;
        }

        static void LogResponse(string path, string body, HttpStatusCode statusCode)
        {
            Debug.Log(path + " " + statusCode.ToString() + " " + body);
        }
    }
}
