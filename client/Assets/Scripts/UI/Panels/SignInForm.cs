using Grpc.Core;
using Idlemon.Data;
using UnityEngine;
using UnityEngine.EventSystems;
using UnityEngine.SceneManagement;
using UnityEngine.UI;

namespace Idlemon.Ui
{
    public class SignInForm : MonoBehaviour
    {
        public InputField email, password;
        public Toggle rememberMe;
        public Button signInBtn;
        public FlashMessage flashMessage;

        Proto.Auth.AuthClient client;

        void Awake()
        {
            client = new Proto.Auth.AuthClient(Grpc.Channel);
        }

        void Start()
        {
            signInBtn.onClick.AddListener(OnBtnClick);

            // Fill in form is remember me data is present
            if (Auth.HasSavedCredentials)
            {
                rememberMe.isOn = true;
                email.text = Auth.SavedEmail;
                password.text = Auth.SavedPassword;
            }
        }

        void Update()
        {
            if (Input.GetKeyDown(KeyCode.Return) || Input.GetKeyDown(KeyCode.KeypadEnter))
            {
                signInBtn.onClick.Invoke();
            }

            if (Input.GetKeyDown(KeyCode.Tab))
            {
                if (EventSystem.current.currentSelectedGameObject == email.gameObject)
                {
                    EventSystem.current.SetSelectedGameObject(password.gameObject);
                }
                else
                {
                    EventSystem.current.SetSelectedGameObject(email.gameObject);
                }
            }
        }

        /// <summary>
        /// Triggered when the sign in button is clicked.
        /// </summary>
        public async void OnBtnClick()
        {
            // validate the form inputs
            if (email.text.Length < 2)
            {
                flashMessage.Flash("username must be longer than 1 character");
                return;
            }
            else if (password.text.Length < 8)
            {
                flashMessage.Flash("password must have at least 8 characters");
                return;
            }

            LoadingPanel.instance.Show();

            try
            {
                var response = await client.SignInAsync(new Proto.SignInRequest { Email = email.text, Pass = password.text }, null, Grpc.Deadline);
                Global.User = new User(response);
                Auth.UpdatePlayerPrefs(email.text, password.text, rememberMe.isOn);
                SceneManager.LoadScene("Overworld");
            }
            catch (RpcException e)
            {
                flashMessage.Flash(e.Status.Detail);
            }
            finally
            {
                LoadingPanel.instance.Hide();
            }
        }
    }
}
