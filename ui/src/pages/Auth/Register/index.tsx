import React, { useState } from 'react';
import { useTranslation } from 'react-i18next';
import { Link } from 'react-router-dom';

import clsx from 'clsx';

import { Button } from '@/components/ui/button';
import {
  Card,
  CardContent,
  CardDescription,
  CardFooter,
  CardHeader,
  CardTitle,
} from '@/components/ui/card';
import { Field, FieldError, FieldLabel } from '@/components/ui/field';
import { Input } from '@/components/ui/input';
import { loginSettingStore } from '@/stores';
import { FormDataType, RegisterReqParams } from '@/common/interface';
import usePageTags from '@/hooks/usePageTags';
import { scrollToElementTop } from '@/utils/common';
import { register } from '@/services/common';

const passwordRules = (password: string) => {
  return {
    minLength: password.length >= 8,
    hasUpper: /[A-Z]/.test(password),
    hasLower: /[a-z]/.test(password),
    hasNumber: /[0-9]/.test(password),
    hasSymbol: /[^A-Za-z0-9]/.test(password),
  };
};

const Index: React.FC = () => {
  const { t } = useTranslation('translation');
  const registerSetting = loginSettingStore((state) => state.login);
  const [formData, setFormData] = useState<FormDataType>({
    email: {
      value: '',
      isInvalid: false,
      errorMsg: '',
    },
    pass: {
      value: '',
      isInvalid: false,
      errorMsg: '',
    },
  });

  const checkValidated = (): boolean => {
    let bol = true;

    const { email, pass } = formData;

    if (!email.value) {
      bol = false;
      formData.email = {
        value: '',
        isInvalid: true,
        errorMsg: t('form.validation.required'),
      };
    }

    if (!pass.value) {
      bol = false;
      formData.pass = {
        value: '',
        isInvalid: true,
        errorMsg: t('form.validation.required'),
      };
    }

    const passRules = passwordRules(pass.value);

    if (!Object.values(passRules).every(Boolean)) {
      bol = false;

      formData.pass = {
        value: pass.value,
        isInvalid: true,
        errorMsg: t('form.validation.password_not_strong_enough'),
      };
    }

    setFormData({
      ...formData,
    });

    if (!bol) {
      const errObj = Object.keys(formData).filter(
        (key) => formData[key].isInvalid,
      );
      const ele = document.getElementById(errObj[0]);
      scrollToElementTop(ele);
    }
    return bol;
  };

  const handleRegister = (event?: any) => {
    if (event) {
      event.preventDefault();
    }

    const params: RegisterReqParams = {
      email: formData.email.value,
      password: formData.pass.value,
    };

    register(params);
  };

  const handleSubmit = (event: any) => {
    event.preventDefault();
    event.stopPropagation();
    if (checkValidated()) {
      handleRegister();
    }
  };

  const rules = passwordRules(formData.pass.value);

  usePageTags({
    title: t('page.sign_up.title'),
  });

  const showFormSignup =
    registerSetting.allow_new_registrations &&
    registerSetting.allow_email_registrations;

  return (
    <Card className="w-full md:w-[460px] md:rounded md:border md:bg-background-card md:shadow-sm">
      <CardHeader>
        <CardTitle className="mb-1">{t('page.sign_up.title')}</CardTitle>
        <CardDescription>{t('page.sign_up.desc')}</CardDescription>
      </CardHeader>
      <CardContent>
        {showFormSignup && (
          <form noValidate onSubmit={handleSubmit}>
            <div className="flex flex-col gap-6">
              <div className="grid">
                <Field className="mb-5">
                  <FieldLabel htmlFor="email">
                    {t('form.label.email')}
                  </FieldLabel>
                  <Input
                    id="email"
                    className="h-11"
                    value={formData.email.value}
                    aria-invalid={formData.email.isInvalid}
                    onChange={(e) =>
                      setFormData((prev) => ({
                        ...prev,
                        email: {
                          ...prev.email,
                          value: e.target.value,
                          isInvalid: false,
                          errorMsg: '',
                        },
                      }))
                    }
                  />
                  <FieldError>{formData.email.errorMsg}</FieldError>
                </Field>
                <Field className="mb-3">
                  <FieldLabel htmlFor="pass">
                    {t('form.label.password')}
                  </FieldLabel>
                  <Input
                    type="password"
                    id="pass"
                    className="h-11"
                    value={formData.pass.value}
                    aria-invalid={formData.pass.isInvalid}
                    onChange={(e) =>
                      setFormData((prev) => ({
                        ...prev,
                        pass: {
                          ...prev.pass,
                          value: e.target.value,
                          isInvalid: false,
                          errorMsg: '',
                        },
                      }))
                    }
                  />
                  <FieldError>{formData.pass.errorMsg}</FieldError>
                  <div className="my-2 space-y-2 text-xs flex justify-center flex-col">
                    <div className="flex items-center">
                      <div
                        className={clsx(
                          'mr-1.5 inline-block h-3 w-3 rounded-full align-middle bg-border',
                          rules.minLength ? 'bg-green-500' : 'bg-muted',
                        )}
                      />
                      <span>
                        {t('form.validation.min_length', { count: 8 })}
                      </span>
                    </div>
                    <div className="flex items-center">
                      <div
                        className={clsx(
                          'mr-1.5 inline-block h-3 w-3 rounded-full align-middle bg-border',
                          rules.hasUpper ? 'bg-green-500' : 'bg-muted',
                        )}
                      />
                      <span>
                        {t('form.validation.has_uppercase', { count: 1 })}
                      </span>
                    </div>
                    <div className="flex items-center">
                      <div
                        className={clsx(
                          'mr-1.5 inline-block h-3 w-3 rounded-full align-middle bg-border',
                          rules.hasLower ? 'bg-green-500' : 'bg-muted',
                        )}
                      />
                      <span>
                        {t('form.validation.has_lowercase', { count: 1 })}
                      </span>
                    </div>
                    <div className="flex items-center">
                      <div
                        className={clsx(
                          'mr-1.5 inline-block h-3 w-3 rounded-full align-middle bg-border',
                          rules.hasNumber ? 'bg-green-500' : 'bg-muted',
                        )}
                      />
                      <span>
                        {t('form.validation.has_number', { count: 1 })}
                      </span>
                    </div>
                    <div className="flex items-center">
                      <div
                        className={clsx(
                          'mr-1.5 inline-block h-3 w-3 rounded-full align-middle bg-border',
                          rules.hasSymbol ? 'bg-green-500' : 'bg-muted',
                        )}
                      />
                      <span>
                        {t('form.validation.has_special_character', {
                          count: 1,
                        })}
                      </span>
                    </div>
                  </div>
                </Field>
                {/* {loginSetting.allow_user_recover && (
                <div className="mb-4 grid justify-items-end text-right">
                  <Link
                    to="/auth/account-recovery"
                    className="text-primary text-sm font-semibold hover:text-primary/80">
                    {t('form.label.forgot_password')}
                  </Link>
                </div>
              )} */}
                <Button type="submit" size="lg">
                  {t('form.sign_up')}
                </Button>
              </div>
            </div>
          </form>
        )}
      </CardContent>
      <CardFooter className="flex flex-col gap-4">
        <div className="w-full mb-3">
          <div className="relative mt-6 mb-4">
            <div className="absolute inset-0 flex items-center">
              <div className="w-full border-t border-border" />
            </div>
            <div className="relative flex justify-center text-sm">
              <span className="bg-background px-3 text-muted-foreground">
                {t('page.or_create_account.title')}
              </span>
            </div>
          </div>
          <section className="relative flex shrink-0 flex-col w-auto text-left">
            <div className="grid gap-3">
              <div className="grid auto-rows-auto grid-cols-3 gap-3">
                <Button variant="outline">
                  <svg
                    width="24px"
                    height="24px"
                    viewBox="0 0 24 24"
                    version="1.1"
                    xmlns="http://www.w3.org/2000/svg">
                    <g
                      id="Artboard"
                      stroke="none"
                      strokeWidth="1"
                      fill="none"
                      fillRule="evenodd">
                      <rect id="Rectangle" x="0" y="0" width="24" height="24" />
                      <g
                        id='Google_"G"_Logo'
                        transform="translate(-0.000000, 0.000000)"
                        fillRule="nonzero">
                        <path
                          d="M23.9775701,12.2813187 C23.9775701,11.4681319 23.9102804,10.6505495 23.766729,9.85054945 L12.251215,9.85054945 L12.251215,14.4571429 L18.8456075,14.4571429 C18.5719626,15.9428571 17.6927103,17.2571429 16.4052336,18.0923077 L16.4052336,21.0813187 L20.3394393,21.0813187 C22.6497196,18.9978022 23.9775701,15.9208791 23.9775701,12.2813187 Z"
                          id="Path"
                          fill="#4285F4"
                        />
                        <path
                          d="M12.251215,23.9692308 C15.5439252,23.9692308 18.3207477,22.9098901 20.3439252,21.0813187 L16.4097196,18.0923077 C15.3151402,18.821978 13.9020561,19.2351648 12.2557009,19.2351648 C9.07065421,19.2351648 6.37009346,17.1296703 5.4011215,14.2989011 L1.34130841,14.2989011 L1.34130841,17.3802198 C3.41383178,21.4197802 7.63514019,23.9692308 12.251215,23.9692308 Z"
                          id="Path"
                          fill="#34A853"
                        />
                        <path
                          d="M5.39663551,14.2989011 C4.88523364,12.8131868 4.88523364,11.2043956 5.39663551,9.71868132 L5.39663551,6.63736264 L1.34130841,6.63736264 C-0.390280374,10.0175824 -0.390280374,14 1.34130841,17.3802198 L5.39663551,14.2989011 Z"
                          id="Path"
                          fill="#FBBC04"
                        />
                        <path
                          d="M12.251215,4.77802198 C13.9917757,4.75164835 15.6740187,5.39340659 16.9345794,6.57142857 L16.9345794,6.57142857 L20.4201869,3.15604396 C18.2130841,1.12527473 15.2837383,0.00879120879 12.251215,0.043956044 C7.63514019,0.043956044 3.41383178,2.59340659 1.34130841,6.63736264 L5.39663551,9.71868132 C6.3611215,6.88351648 9.06616822,4.77802198 12.251215,4.77802198 Z"
                          id="Path"
                          fill="#EA4335"
                        />
                      </g>
                    </g>
                  </svg>
                  Google
                </Button>
              </div>
            </div>
          </section>
        </div>
        <div className="text-center w-full">
          <span className="text-foreground text-sm">
            {t('page.have_account.title')}
          </span>
          <Link
            to="/auth/login"
            className="ml-1 text-primary text-sm font-semibold hover:text-primary/80">
            {t('form.sign_in')}
          </Link>
        </div>
      </CardFooter>
    </Card>
  );
};

export default React.memo(Index);
