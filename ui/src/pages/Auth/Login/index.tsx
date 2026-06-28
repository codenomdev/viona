import React, { FormEvent, useState } from 'react';
import { Link } from 'react-router-dom';

import usePageTags from '@/hooks/usePageTags';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Field, FieldError, FieldLabel } from '@/components/ui/field';
import { Input } from '@/components/ui/input';
import { Button } from '@/components/ui/button';
import { FormDataType, LoginReqParams } from '@/common/interface';
import { handleFormError, scrollToElementTop } from '@/utils';
import { login } from '@/services';

const Index: React.FC = () => {
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
        errorMsg: 'required',
      };
    }

    if (!pass.value) {
      bol = false;
      formData.pass = {
        value: '',
        isInvalid: true,
        errorMsg: 'required',
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
    title: 'Login Users',
  });
  return (
    <Card className="w-full md:w-[460px] md:rounded md:border md:bg-background-card md:shadow-sm">
      <CardHeader>
        <CardTitle>Sign in</CardTitle>
      </CardHeader>
      <CardContent>
        <form noValidate onSubmit={handleSubmit}>
          <div className="flex flex-col gap-6">
            <div className="grid">
              <Field className="mb-5">
                <FieldLabel htmlFor="email">Email</FieldLabel>
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
                <FieldLabel htmlFor="pass">Password</FieldLabel>
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
              <div className="mb-4 grid justify-items-end text-right">
                <Link
                  to="auth/account-recovery"
                  className="text-primary text-sm font-semibold hover:text-primary/80">
                  Forgot password?
                </Link>
              </div>
              <Button type="submit" size="lg">
                Sign in
              </Button>
            </div>
          </div>
        </form>
      </CardContent>
    </Card>
  );
};

export default React.memo(Index);
