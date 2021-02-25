using UnityEngine;
using UnityEngine.UI;

namespace Idlemon.Ui
{
    public class RefreshRoomsListBtn : MonoBehaviour
    {
        public LobbyListPanel lobbyListPanel;

        void Start()
        {
            var btn = GetComponent<Button>();
            btn.onClick.AddListener(OnClick);
        }

        async void OnClick()
        {
            await lobbyListPanel.RefreshRoomsList();
        }
    }
}
