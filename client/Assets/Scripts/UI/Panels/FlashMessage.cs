using UnityEngine;
using UnityEngine.UI;
using System.Collections;

namespace Idlemon.Ui
{
    /// <summary>
    /// Used to flash messages to the player.
    /// </summary>
    public class FlashMessage : MonoBehaviour
    {
        public float lifetime = 2.5f;
        public Text text;

        public void Flash(string message)
        {
            text.text = message;
            Show();

            StopAllCoroutines();
            StartCoroutine(R_FadeOut());
        }

        void Show()
        {
            GetComponent<Image>().enabled = true;
            text.enabled = true;
        }

        public void Hide()
        {
            GetComponent<Image>().enabled = false;
            text.enabled = false;
        }

        IEnumerator R_FadeOut()
        {
            yield return new WaitForSeconds(lifetime);
            Hide();
        }
    }
}
