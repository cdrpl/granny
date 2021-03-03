using Idlemon.Data;
using UnityEngine;
using UnityEngine.EventSystems;
using UnityEngine.UI;

namespace Idlemon.Ui
{
    public class ChatPanel : MonoBehaviour
    {
        public InputField input;
        public Transform scrollContent;
        public GameObject scrollItem; // prefab for message scroll items

        void Start()
        {
            /*webSocketClient.OnChat.AddListener((msg) =>
            {
                var chatMessage = new ChatMessage(msg);
                AttachMessage(chatMessage);
            });*/
        }

        void Update()
        {
            if (Input.GetKeyDown(KeyCode.Return) || Input.GetKeyDown(KeyCode.KeypadEnter))
            {
                GameObject current = EventSystem.current.currentSelectedGameObject;

                if (current == input.gameObject)
                {
                    SendMessage();

#if !UNITY_ANDROID
                    input.ActivateInputField(); // re-select the input field after sending a message
#endif
                }
            }
        }

        public void SendMessage()
        {
            if (input.text != string.Empty)
            {
                /*var chatMessage = new ChatMessage(player.name, input.text);

                AttachMessage(chatMessage);

                // send WebSocket message
                byte[] bytes = System.Text.Encoding.ASCII.GetBytes(input.text);
                var msg = new SocketMessage(Const.Channel.Chat, bytes);
                webSocketClient.SendWebSocketMessage(msg);

                input.text = string.Empty;*/
            }
        }

        /// <summary>
        /// Attach the message to the UI.
        /// </summary>
        /*public void AttachMessage(ChatMessage chatMessage)
        {
            // instantiate item
            Transform item = Instantiate<GameObject>(scrollItem).transform;
            item.SetParent(scrollContent);
            item.localScale = Vector3.one;

            // update item ui
            item.GetComponent<ChatMessageScrollItem>().Draw(chatMessage);
        }*/
    }
}
