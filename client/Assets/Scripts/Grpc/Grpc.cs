using Grpc.Core;
using System;
using System.Threading;

namespace Idlemon
{
    public static class Grpc
    {
        /// <summary>
        /// Grpc timeout in seconds.
        /// </summary>
        public const int TIMEOUT = 5;

        /// <summary>
        /// Single gRPC channel used throughout the app.
        /// </summary>
        public static readonly Channel Channel = new Channel(Const.SERVER_ADDR, ChannelCredentials.Insecure);

        /// <summary>
        /// Default cancellation token source for timeouts.
        /// </summary>
        public static readonly CancellationTokenSource Cts = new CancellationTokenSource(TimeSpan.FromSeconds(5));

        /// <summary>
        /// The default gRPC deadline.
        /// </summary>
        /// <returns></returns>
        public static DateTime Deadline => DateTime.UtcNow.AddSeconds(TIMEOUT);

        /// <summary>
        /// Metadata holding the user ID and auth token.
        /// </summary>
        public static Metadata Metadata
        {
            get
            {
                var meta = new Metadata();
                meta.Add("user-id", Global.User.Id.ToString());
                meta.Add("token", Global.User.Token);
                return meta;
            }
        }
    }
}
