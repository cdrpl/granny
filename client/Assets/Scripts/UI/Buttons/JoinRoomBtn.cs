using UnityEngine;
using UnityEngine.UI;

namespace Idlemon.Ui
{
    public class JoinRoomBtn : MonoBehaviour
    {
        public JoinRoom joinRoom;

        Button button;

        void Awake()
        {
            button = GetComponent<Button>();
        }

        void Start()
        {
            button.onClick.AddListener(OnBtnPress);
        }

        async void OnBtnPress()
        {
            await joinRoom.Join();
        }
    }
}
