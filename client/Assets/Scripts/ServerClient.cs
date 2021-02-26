using NativeWebSocket;
using System.Collections.Generic;
using UnityEngine;

namespace Idlemon
{
    public class ServerClient : MonoBehaviour
    {
        WebSocket websocket;

        public async void Connect()
        {
            if (Global.User == null)
            {
                Debug.LogWarning("Global user is null, LobbyClient will not connect");
                return;
            }

            // Put authorization details in header
            string authorization = Global.User.Id.ToString() + ":" + Global.User.Token;
            var header = new Dictionary<string, string>
            {
                {"authorization", authorization},
            };

            websocket = new WebSocket("ws://" + Const.SERVER_ADDR + "/ws", header);

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
                string message = System.Text.Encoding.UTF8.GetString(bytes);
                Debug.Log("Recv message: " + message);
            };

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

        async void SendWebSocketMessage(string message)
        {
            if (websocket.State == WebSocketState.Open)
            {
                await websocket.SendText(message);
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
