using Grpc.Core;
using Idlemon.Ui;
using Proto;
using System.Threading.Tasks;
using UnityEngine;

namespace Idlemon
{
    /// <summary>
    /// User joined server stream.
    /// </summary>
    public class UserJoined
    {
        public FlashMessage flashMessage;

        Proto.Room.RoomClient client;

        public UserJoined()
        {
            client = new Proto.Room.RoomClient(Grpc.Channel);
        }

        public async Task Stream()
        {
            try
            {
                // Setup the join room stream
                using var stream = client.UserJoined(new UserJoinedReq(), Grpc.Metadata, Grpc.Deadline);
                while (await stream.ResponseStream.MoveNext())
                {
                    Debug.Log(stream.ResponseStream.Current.Id);
                }
                Debug.Log("Stream ended");
            }
            catch (RpcException e)
            {
                flashMessage.Flash(e.Status.Detail);
            }
        }
    }
}
