using NativeWebSocket;
using System.Collections.Generic;
using UnityEngine;

namespace Idlemon
{
    public class LobbyClient : MonoBehaviour
    {
        WebSocket websocket;

        public async void Connect()
        {
            if (Global.User == null)
            {
                Debug.LogWarning("Global user is null, LobbyClient will not connect");
                return;
            }

            // Authentication details in header
            var header = new Dictionary<string, string>
            {
                {"id", Global.User.Id.ToString()},
                {"token", Global.User.Token},
            };

            websocket = new WebSocket("ws://" + Const.LOBBY_ADDR + "/ws", header);

            websocket.OnOpen += () =>
            {
                Debug.Log("Connection open!");
            };

            websocket.OnError += (e) =>
            {
                Debug.Log("Error! " + e);
            };

            websocket.OnClose += (e) =>
            {
                Debug.Log("WebSocketCloseCode:" + e.ToString());
            };

            websocket.OnMessage += (bytes) =>
            {
                Debug.Log("OnMessage!");
                Debug.Log(bytes);

                // getting the message as a string
                // var message = System.Text.Encoding.UTF8.GetString(bytes);
                // Debug.Log("OnMessage! " + message);
            };

            // Keep sending messages at every 0.3s
            InvokeRepeating("SendWebSocketMessage", 0.0f, 0.3f);

            // waiting for messages
            await websocket.Connect();
        }

        void Update()
        {
            if (Input.GetKeyDown(KeyCode.C))
            {
                Connect();
            }

#if !UNITY_WEBGL || UNITY_EDITOR
            if (websocket != null)
            {
                websocket.DispatchMessageQueue();
            }
#endif
        }

        async void SendWebSocketMessage()
        {
            if (websocket.State == WebSocketState.Open)
            {
                await websocket.SendText("Client says hello");
            }
        }

        async void OnDestroy()
        {
            if (websocket != null)
            {
                await websocket.Close();
            }
        }

    }
}
