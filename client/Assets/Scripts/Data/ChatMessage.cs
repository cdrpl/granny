using System;

namespace Idlemon.Data
{
    public class ChatMessage
    {
        /// <summary>
        /// The user ID of the message sender.
        /// </summary>
        public string UserId { get; set; }

        /// <summary>
        /// The name of the message sender.
        /// </summary>
        public string Name { get; set; }

        /// <summary>
        /// The chat message.
        /// </summary>
        public string Message { get; set; }

        /// <summary>
        /// The time the message was sent at.
        /// </summary>
        public DateTime Time { get; set; }
    }
}
