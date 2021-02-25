using Idlemon.Data;
using UnityEngine;
using UnityEngine.UI;

namespace Idlemon.Ui
{
    public class UsernameLabel : MonoBehaviour
    {
        public User UserData
        {
            get => userData;
            set
            {
                userData = value;
            }
        }

        Text label;
        User userData;

        void Awake()
        {
            label = GetComponent<Text>();
        }

        void Start()
        {
            userData = Global.User;
            UpdateLabel();
        }

        void UpdateLabel()
        {
            label.text = userData.Name;
        }
    }
}
