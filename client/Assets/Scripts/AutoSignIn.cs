using Grpc.Core;
using Idlemon.Data;
using Idlemon.Ui;
using UnityEngine;

namespace Idlemon
{
    /// <summary>
    /// Will automatically authenticate the user using details saved in player prefs.
    /// </summary>
    public class AutoSignIn : MonoBehaviour
    {
        public FlashMessage flashMessage;

        Proto.Auth.AuthClient client;

        async void Awake()
        {
            if (Global.User != null)
            {
                return;
            }

            client = new Proto.Auth.AuthClient(Grpc.Channel);

            try
            {
                LoadingPanel.instance.Show();

                var req = new Proto.SignInRequest { Email = Auth.SavedEmail, Pass = Auth.SavedPassword };
                var response = await client.SignInAsync(req, null, Grpc.Deadline);
                Global.User = new User(response);
                Debug.Log("User logged in: " + Auth.SavedEmail);
            }
            catch (RpcException e)
            {
                flashMessage.Flash(e.Status.Detail);
            }
            finally
            {
                LoadingPanel.instance.Hide();
            }
        }
    }
}
