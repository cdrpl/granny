using Idlemon.Data;
using UnityEngine;
using UnityEngine.UI;

namespace Idlemon.Ui
{
    /// <summary>
    /// The scroll item used for displaying the hosted rooms in the lobby.
    /// </summary>
    public class RoomScrollItem : MonoBehaviour
    {
        public Text nameText, playersText;
        public Button joinBtn;

        public Room Room { get; set; }

        void Start()
        {
            joinBtn.onClick.AddListener(OnBtnClick);
        }

        /// <summary>
        /// Updates the UI elements to represent the assigned room.
        /// </summary>
        public void Redraw()
        {
            if (Room == null)
            {
                Debug.LogError("Redraw attempt while Room is null", this);
                return;
            }

            nameText.text = Room.Name;
            playersText.text = Room.NumUsers.ToString() + "/5";
        }

        void OnBtnClick()
        {
            if (Room == null)
            {
                Debug.LogError("Attempt to join room but Room is null", this);
                return;
            }

            Debug.Log("Join room");
        }
    }
}
