using Grpc.Core;
using Idlemon.Ui;
using System.Threading.Tasks;
using UnityEngine;

namespace Idlemon
{
    public class JoinRoom : MonoBehaviour
    {
        public FlashMessage flashMessage;

        Proto.Room.RoomClient client;

        void Awake()
        {
            client = new Proto.Room.RoomClient(Grpc.Channel);
        }

        public async Task<bool> Join()
        {
            try
            {
                LoadingPanel.instance.Show();
                await client.JoinRoomAsync(new Proto.JoinRoomReq(), Grpc.Metadata, Grpc.Deadline);
                Debug.Log("Joined room");
                return true;
            }
            catch (RpcException e)
            {
                flashMessage.Flash(e.Status.Detail);
            }
            finally
            {
                LoadingPanel.instance.Hide();
            }

            return false;
        }
    }
}
