using Grpc.Core;
using UnityEngine;

namespace Idlemon
{
    public class GetRoom : MonoBehaviour
    {
        Channel channel;
        Proto.Room.RoomClient client;

        void Awake()
        {
            channel = new Channel(Const.SERVER_ADDR, ChannelCredentials.Insecure);
            client = new Proto.Room.RoomClient(channel);
        }

        async void Rpc()
        {
            var room = await client.GetRoomAsync(new Proto.GetRoomRequest { });
        }
    }
}
