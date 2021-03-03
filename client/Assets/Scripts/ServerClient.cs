using NativeWebSocket;
using System.Collections.Generic;
using UnityEngine;

namespace Idlemon
{
    public class ServerClient : MonoBehaviour
    {
        enum Channel
        {
            JoinRoom
        }

        WebSocket websocket;

        void Start()
        {
            if (Global.User != null)
            {
                Connect();
            }
        }

        public async void Connect()
        {
            // Put authorization details in header
            string authorization = Global.User.Id.ToString() + ":" + Global.User.Token;
            var header = new Dictionary<string, string>
            {
                {"authorization", authorization},
            };

            websocket = new WebSocket("ws://" + Const.SERVER_ADDR + "/ws", header);

            websocket.OnOpen += () =>
            {
                Debug.Log("Server client connected");
            };

            websocket.OnError += (e) =>
            {
                Debug.Log("Server client error: " + e);
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

        async void Update()
        {
            if (Input.GetKeyDown(KeyCode.Tab))
            {
                await websocket.Send(new byte[] { (byte)Channel.JoinRoom });
            }
            if (Input.GetKeyDown(KeyCode.R))
            {
                string res = await Web.GetRoom();
                Debug.Log(res);
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
