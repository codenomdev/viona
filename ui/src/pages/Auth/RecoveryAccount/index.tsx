import React, { useState } from 'react';
import { useTranslation } from 'react-i18next';
import { useNavigate } from 'react-router-dom';

// import { Link } from 'lucide-react';

import { ChevronLeft } from 'lucide-react';

import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from '@/components/ui/card';
import { Field, FieldError, FieldLabel } from '@/components/ui/field';
import { Input } from '@/components/ui/input';
import { Button } from '@/components/ui/button';
import { FormDataType, RecoveryAccountReqParams } from '@/common/interface';
import { scrollToElementTop } from '@/utils/common';
import { recoveryAccount } from '@/services/common';
import usePageTags from '@/hooks/usePageTags';
import { loginSettingStore } from '@/stores';

const Index: React.FC = () => {
  const { t } = useTranslation('translation');
  const recoverySetting = loginSettingStore((state) => state.login);
  const [formData, setFormData] = useState<FormDataType>({
    email: {
      value: '',
      isInvalid: false,
      errorMsg: '',
    },
  });

  const navigate = useNavigate();

  const checkValidated = (): boolean => {
    let bol = true;

    const { email } = formData;

    if (!email.value) {
      bol = false;
      formData.email = {
        value: '',
        isInvalid: true,
        errorMsg: t('form.validation.required'),
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

  const handleRecovery = (event?: any) => {
    if (event) {
      event.preventDefault();
    }

    const param: RecoveryAccountReqParams = {
      email: formData.email.value,
    };
    recoveryAccount(param);
  };

  const handleSubmit = (event: any) => {
    event.preventDefault();
    event.stopPropagation();
    if (checkValidated()) {
      handleRecovery();
    }
  };

  const showFormRecover = recoverySetting.allow_user_recover;

  usePageTags({
    title: t('page.account_recovery.title'),
  });

  return (
    <>
      <Button
        variant="link"
        className="p-0 gap-0 mb-1 text-xs"
        size="sm"
        onClick={() => navigate('/auth/login')}>
        <ChevronLeft />
        {t('page.back_to_sign_in.title')}
      </Button>
      <Card className="w-full md:w-[460px] md:rounded md:border md:bg-background-card md:shadow-sm">
        <CardHeader>
          <CardTitle className="mb-2">
            {t('page.account_recovery.title')}
          </CardTitle>
          <CardDescription>{t('page.account_recovery.desc')}</CardDescription>
        </CardHeader>
        <CardContent>
          {showFormRecover && (
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
                  <Button type="submit" size="lg">
                    {t('form.send_recovery_email')}
                  </Button>
                </div>
              </div>
            </form>
          )}
        </CardContent>
      </Card>
    </>
  );
};

export default React.memo(Index);
