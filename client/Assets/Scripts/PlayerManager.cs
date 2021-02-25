using UnityEngine;
using System;
using System.Collections.Generic;

namespace Idlemon
{
    /// <summary>
    /// Player manager will track the players and update their positions.
    /// </summary>
    public class PlayerManager : MonoBehaviour
    {
        public GameObject playerPrefab;

        Dictionary<uint, Transform> players = new Dictionary<uint, Transform>();

        void Start()
        {
            //webSocketClient.OnPosition.AddListener(OnPositionReceived);
            //webSocketClient.OnMessage.AddListener(OnSocketMessage);
        }

        /*void OnPositionReceived(SocketMessage message)
        {
            uint id = BitConverter.ToUInt32(message.data, 0);
            float x = BitConverter.ToSingle(message.data, 4);
            float y = BitConverter.ToSingle(message.data, 8);

            if (players.ContainsKey(id))
            {
                players[id].transform.position = new Vector3(x, y, 0);
            }
            else
            {
                SpawnPlayer(id, x, y);
            }
        }

        void OnSocketMessage(SocketMessage message)
        {
            switch (message.channel)
            {
                case (int)Const.Channel.PlayerConnected:
                    uint id = BitConverter.ToUInt32(message.data, 0);

                    if (!players.ContainsKey(id))
                    {
                        float x = BitConverter.ToSingle(message.data, 4);
                        float y = BitConverter.ToSingle(message.data, 8);
                        SpawnPlayer(id, x, y);

                        Debug.Log("Player connected: " + id);
                    }

                    break;

                case (int)Const.Channel.PlayerDisconnected:
                    id = BitConverter.ToUInt32(message.data, 0);
                    Debug.Log("Player disconnected: " + id);

                    if (players.ContainsKey(id))
                    {
                        // Despawn player
                        Destroy(players[id].gameObject);

                        // Remove player from dictionary
                        players.Remove(id);
                    }

                    break;
            }
        }*/

        void SpawnPlayer(uint id, float x, float y)
        {
            // Spawn player
            var pos = new Vector3(x, y, 0);
            var instance = Instantiate<GameObject>(playerPrefab, pos, Quaternion.identity);

            // Add player to dictionary
            players[id] = instance.transform;
        }
    }
}
