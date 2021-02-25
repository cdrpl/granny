using UnityEngine;

namespace Idlemon
{
    public class QuitApp : MonoBehaviour
    {
        public void Quit()
        {
            Application.Quit(0);
        }
    }
}
