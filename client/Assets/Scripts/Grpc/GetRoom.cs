using Grpc.Core;
using UnityEngine;

namespace Idlemon
{
    public class GetRoom : MonoBehaviour
    {
        Proto.Room.RoomClient client;

        void Awake()
        {
            client = new Proto.Room.RoomClient(Grpc.Channel);
        }

        async void Rpc()
        {
            var room = await client.GetRoomAsync(new Proto.GetRoomRequest { });
        }
    }
}
