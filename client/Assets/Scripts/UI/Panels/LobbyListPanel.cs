using Idlemon.Data;
using Newtonsoft.Json;
using System;
using System.Collections.Generic;
using System.Linq;
using System.Threading.Tasks;
using UnityEngine;


namespace Idlemon.Ui
{
    /// <summary>
    /// Handles the panel that displays the list of lobby rooms.
    /// </summary>
    public class LobbyListPanel : MonoBehaviour
    {
        public Transform content;
        public GameObject scrollItemPrefab;
        public HttpLoadingPanel loadingPanel;

        /// <summary>
        /// List of instantiated scroll items.
        /// </summary>
        List<RoomScrollItem> scrollItems;

        void Awake()
        {
            scrollItems = new List<RoomScrollItem>();
        }

        async void Start()
        {
            await RefreshRoomsList();
        }

        public async Task RefreshRoomsList()
        {
            try
            {
                loadingPanel.Show();

                ApiResponse response = await Web.GetRooms();

                // Destroy currently spawned scroll items
                foreach (RoomScrollItem item in scrollItems.ToList())
                {
                    scrollItems.Remove(item);
                    Destroy(item.gameObject);
                }

                // Deserialize the API response
                Room[] rooms = JsonConvert.DeserializeObject<Room[]>(response.Success.Data.ToString());

                // Spawn the scroll items
                foreach (Room room in rooms)
                {
                    var go = Instantiate<GameObject>(scrollItemPrefab);
                    go.transform.SetParent(content);
                    go.transform.localScale = Vector3.one;

                    var item = go.GetComponent<RoomScrollItem>();
                    scrollItems.Add(item);

                    item.Room = room;
                    item.Redraw();
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
