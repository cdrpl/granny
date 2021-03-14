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
            bool joined = await joinRoom.Join();

            if (joined)
            {
                var userJoined = new UserJoined();
                await userJoined.Stream();
                Debug.Log("end");
            }
        }
    }
}
