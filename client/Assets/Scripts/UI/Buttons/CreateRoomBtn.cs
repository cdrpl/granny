using Idlemon.Data;
using Newtonsoft.Json;
using System;
using System.Net;
using UnityEngine;
using UnityEngine.UI;

namespace Idlemon.Ui
{
    public class CreateRoomBtn : MonoBehaviour
    {
        public Button button;
        public HttpLoadingPanel loadingPanel;
        public FlashMessage flashMessage;

        void Awake()
        {
            button.onClick.AddListener(OnClick);
        }

        async void OnClick()
        {
            try
            {
                loadingPanel.Show();

                ApiResponse response = await Web.CreateRoom(Global.User.Id, "Noobs welcome");

                if (response.StatusCode == HttpStatusCode.OK)
                {
                    Debug.Log(response.Success.Data);
                    string data = response.Success.Data.ToString();
                    Room room = JsonConvert.DeserializeObject<Room>(data);
                }
                else if (response.HasError)
                {
                    flashMessage.Flash(response.Error.Message);
                }
            }
            catch (Exception e)
            {
                Debug.LogError(e, this);
            }
            finally
            {
                loadingPanel.Hide();
            }
        }
    }
}
