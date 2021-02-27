using System.Net;
using UnityEngine;

namespace Idlemon
{
    /// <summary>
    /// Will automatically authenticate the user using details saved in player prefs.
    /// </summary>
    public class AutoSignIn : MonoBehaviour
    {
        public ServerClient serverClient;

        async void Awake()
        {
            if (Global.User != null)
            {
                return;
            }

            var response = await Auth.SignIn(Auth.SavedEmail, Auth.SavedPassword, true);

            if (response.StatusCode == HttpStatusCode.OK)
            {
                Debug.Log("User has logged on: " + Global.User.Name);
                serverClient.Connect();
            }
            else
            {
                Debug.LogWarning("User login failed: " + response.Error.Message);
            }
        }
    }
}
