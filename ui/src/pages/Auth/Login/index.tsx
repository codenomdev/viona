import React, { FormEvent, useState } from 'react';
import { Link } from 'react-router-dom';
import { useTranslation } from 'react-i18next';

import usePageTags from '@/hooks/usePageTags';
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
import { Button } from '@/components/ui/button';
import { FormDataType, LoginReqParams } from '@/common/interface';
import { handleFormError, scrollToElementTop } from '@/utils';
import { login } from '@/services';
import { loginSettingStore } from '@/stores';

const Index: React.FC = () => {
  const { t } = useTranslation('translation');
  const loginSetting = loginSettingStore((state) => state.login);
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

    setFormData({
      ...formData,
    });

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

  const handleLogin = (event?: any) => {
    if (event) {
      event.preventDefault();
    }
    const params: LoginReqParams = {
      email: formData.email.value,
      password: formData.pass.value,
    };

    login(params)
      .then(async (res) => {
        console.log('response: ', res);
      })
      .catch((err) => {
        if (err.isError) {
          const data = handleFormError(err, formData);
          setFormData({ ...data });
          // passwordCaptcha?.handleCaptchaError?.(err.list);
          const ele = document.getElementById(err.list[0].error_field);
          scrollToElementTop(ele);
        }
      });
  };

  const handleSubmit = async (event: FormEvent) => {
    event.preventDefault();
    event.stopPropagation();

    if (!checkValidated()) {
      return;
    }

    // if (!passwordCaptcha) {
    //   handleLogin();
    //   return;
    // }

    // passwordCaptcha?.check?.(() => {
    //   handleLogin();
    // });
    handleLogin();
  };

  usePageTags({
    title: t('page.sign_in.title'),
  });
  return (
    <Card className="w-full md:w-[460px] md:rounded md:border md:bg-background-card md:shadow-sm">
      <CardHeader>
        <CardTitle className="mb-1">{t('page.sign_in.title')}</CardTitle>
        <CardDescription>{t('page.sign_in.desc')}</CardDescription>
      </CardHeader>
      <CardContent>
        <form noValidate onSubmit={handleSubmit}>
          <div className="flex flex-col gap-6">
            <div className="grid">
              <Field className="mb-5">
                <FieldLabel htmlFor="email">{t('form.label.email')}</FieldLabel>
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
              </Field>
              {loginSetting.allow_user_recover && (
                <div className="mb-4 grid justify-items-end text-right">
                  <Link
                    to="/auth/account-recovery"
                    className="text-primary text-sm font-semibold hover:text-primary/80">
                    {t('form.label.forgot_password')}
                  </Link>
                </div>
              )}
              <Button type="submit" size="lg">
                {t('form.sign_in')}
              </Button>
            </div>
          </div>
        </form>
      </CardContent>
      <CardFooter>
        <div className="text-center w-full">
          <span className="text-foreground text-sm">
            {t('page.not_have_account.title')}
          </span>
          <Link
            to="/auth/register"
            className="ml-1 text-primary text-sm font-semibold hover:text-primary/80">
            {t('form.sign_up')}
          </Link>
        </div>
      </CardFooter>
    </Card>
  );
};

export default React.memo(Index);
