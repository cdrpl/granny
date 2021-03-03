using UnityEngine;
using UnityEngine.UI;
using System.Collections;

namespace Idlemon.Ui
{
    /// <summary>
    /// Controls the panel that appears when a network request is waiting for a response.
    /// </summary>
    public class LoadingPanel : MonoBehaviour
    {
        public static LoadingPanel instance { get; private set; }

        public float loadingImageDelay = 0.75f; // delay in seconds before showing the loading image
        public float timeout = 10f; // disable the panel if timeout has been reached

        void Awake()
        {
            if (instance == null)
            {
                instance = this;
            }
        }

        public void Show()
        {
            StartCoroutine(ShowLoadingImage(loadingImageDelay)); // show loading image
        }

        public void Hide()
        {
            StopAllCoroutines();

            GetComponent<Image>().enabled = false;
            transform.GetChild(0).gameObject.SetActive(false);
        }

        private IEnumerator ShowLoadingImage(float delay)
        {
            yield return new WaitForSeconds(delay);

            // image covers the whole screen to block inputs
            GetComponent<Image>().enabled = true;

            transform.GetChild(0).gameObject.SetActive(true); // enable loading animation

            // timeout handling
            yield return new WaitForSeconds(timeout);
            Hide();
        }
    }
}
