using Grpc.Core;

namespace Idlemon
{
    public static class Grpc
    {
        /// <summary>
        /// Single gRPC channel used throughout the app.
        /// </summary>
        public static readonly Channel Channel = new Channel(Const.SERVER_ADDR, ChannelCredentials.Insecure);
    }
}
