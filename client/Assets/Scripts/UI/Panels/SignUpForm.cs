using System;
using System.Net;
using UnityEngine;
using UnityEngine.Events;
using UnityEngine.EventSystems;
using UnityEngine.UI;

namespace Idlemon.Ui
{
    public class SignUpForm : MonoBehaviour
    {
        public InputField username, email, password, password2;
        public Button signUpBtn;
        public FlashMessage flashMessage;
        public HttpLoadingPanel loadingPanel;

        /// <summary>
        /// Triggered on user sign up success.
        /// </summary>
        public UnityEvent OnSignUpSuccess;

        void Start()
        {
            signUpBtn.onClick.AddListener(OnBtnClick);
        }

        void Update()
        {
            if (Input.GetKeyDown(KeyCode.Return) || Input.GetKeyDown(KeyCode.KeypadEnter))
            {
                signUpBtn.onClick.Invoke();
            }

            if (Input.GetKeyDown(KeyCode.Tab))
            {
                if (EventSystem.current.currentSelectedGameObject == username.gameObject)
                {
                    SelectInput(email);
                }
                else if (EventSystem.current.currentSelectedGameObject == email.gameObject)
                {
                    SelectInput(password);
                }
                else if (EventSystem.current.currentSelectedGameObject == password.gameObject)
                {
                    SelectInput(password2);
                }
                else
                {
                    SelectInput(username);
                }
            }
        }

        void OnEnable()
        {
            username.Select();
        }

        async void OnBtnClick()
        {
            // validate input fields
            if (username.text == string.Empty) // username is required
            {
                SelectInput(username);
                flashMessage.Flash("username is required");
                return;
            }
            else if (email.text == string.Empty) // email is required
            {
                SelectInput(email);
                flashMessage.Flash("email is required");
                return;
            }
            else if (password.text == string.Empty) // password is required
            {
                SelectInput(password);
                flashMessage.Flash("password is required");
                return;
            }
            else if (password.text.Length < 8) // password must be 8 chars long
            {
                SelectInput(password);
                flashMessage.Flash("password must contain 8 characters");
                return;
            }
            else if (password.text != password2.text) // passwords must match
            {
                SelectInput(password2);
                flashMessage.Flash("passwords do not match");
                return;
            }

            // Send http request
            try
            {
                loadingPanel.Show();

                var response = await Web.SignUp(username.text, email.text, password.text);

                if (response.HasError)
                {
                    flashMessage.Flash(response.Error.Message);
                }
                else
                {
                    flashMessage.Flash("Account created!");
                    ClearInputs();
                    OnSignUpSuccess.Invoke();
                }
            }
            catch (Exception e)
            {
                Debug.LogError(e, this);
                flashMessage.Flash("An error has occured, check your internet connection");
            }
            finally
            {
                loadingPanel.Hide();
            }
        }

        /// <summary>
        /// Select the specified input.
        /// </summary>
        void SelectInput(InputField input)
        {
            input.ActivateInputField();
            input.Select();
        }

        /// <summary>
        /// Clears the input fields of the form.
        /// </summary>
        void ClearInputs()
        {
            username.text = string.Empty;
            email.text = string.Empty;
            password.text = string.Empty;
            password2.text = string.Empty;
        }
    }
}
