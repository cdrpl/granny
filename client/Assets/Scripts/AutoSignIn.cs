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

        Channel channel;
        Proto.Auth.AuthClient client;

        async void Awake()
        {
            if (Global.User != null)
            {
                return;
            }

            channel = new Channel(Const.SERVER_ADDR, ChannelCredentials.Insecure);
            client = new Proto.Auth.AuthClient(channel);

            try
            {
                LoadingPanel.instance.Show();
                var response = await client.SignInAsync(new Proto.SignInRequest { Email = Auth.SavedEmail, Pass = Auth.SavedPassword });
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
