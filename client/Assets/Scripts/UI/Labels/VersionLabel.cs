using UnityEngine;
using UnityEngine.UI;

namespace Idlemon.Ui
{
    public class VersionLabel : MonoBehaviour
    {
        void Awake()
        {
            GetComponent<Text>().text = "v" + Application.version;
        }
    }
}
