using Idlemon.Data;
using Newtonsoft.Json;
using System.Net;
using System.Threading.Tasks;
using UnityEngine;

namespace Idlemon
{
    public class Auth
    {
        const string EMAIL_PREF_KEY = "email";
        const string PASS_PREF_KEY = "pass";

        public static string SavedEmail => PlayerPrefs.GetString(EMAIL_PREF_KEY);
        public static string SavedPassword => PlayerPrefs.GetString(PASS_PREF_KEY);
        public static bool HasSavedCredentials => PlayerPrefs.HasKey(EMAIL_PREF_KEY);

        /// <summary>
        /// Will return true if successfully signed in.
        /// </summary>
        /// <returns></returns>
        public async static Task<ApiResponse> SignIn(string email, string password, bool rememberMe)
        {
            var response = await Web.SignIn(email, password);

            if (response.StatusCode == HttpStatusCode.OK)
            {
                Global.User = JsonConvert.DeserializeObject<User>(response.Success.Data.ToString());

                UpdatePlayerPrefs(email, password, rememberMe);

                // Set the Authorization header for the HTTP Client.
                string authorization = Global.User.Id.ToString() + ":" + Global.User.Token;
                Web.Client.DefaultRequestHeaders.TryAddWithoutValidation("Authorization", authorization);
            }

            return response;
        }

        static void UpdatePlayerPrefs(string email, string password, bool rememberMe)
        {
            if (rememberMe)
            {
                PlayerPrefs.SetString(EMAIL_PREF_KEY, email);
                PlayerPrefs.SetString(PASS_PREF_KEY, password);
            }
            else
            {
                PlayerPrefs.DeleteKey(EMAIL_PREF_KEY);
                PlayerPrefs.DeleteKey(PASS_PREF_KEY);
            }
        }
    }
}
