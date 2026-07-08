import type * as Type from '@/common/interface';

function scrollToElementTop(element) {
  if (!element) {
    return;
  }
  const offset = 120;
  const bodyRect = document.body.getBoundingClientRect().top;
  const elementRect = element.getBoundingClientRect().top;
  const elementPosition = elementRect - bodyRect;
  const offsetPosition = elementPosition - offset;

  window.scrollTo({
    top: offsetPosition,
    behavior: 'instant' as ScrollBehavior,
  });
}

function handleFormError(
  errorList: Type.FieldError[],
  data: any,
  keymap?: Array<{ from: string; to: string }>,
) {
  if (errorList?.length > 0) {
    errorList.forEach((item) => {
      if (keymap?.length) {
        const key = keymap.find((k) => k.from === item.error_field);

        if (key) {
          item.error_field = key.to;
        }
      }

      const errorFieldObject = data[item.error_field];

      if (errorFieldObject) {
        errorFieldObject.isInvalid = true;
        errorFieldObject.errorMsg = item.error_msg;
      }
    });
  }

  return data;
}

export { scrollToElementTop, handleFormError };
