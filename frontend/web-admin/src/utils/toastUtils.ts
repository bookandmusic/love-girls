import { inject } from "vue";

export const useToast = () => {
  const showToast =
    inject<(message: string, type?: "success" | "error" | "info") => void>(
      "showToast",
    );

  if (!showToast) {
    throw new Error(
      "useToast must be used within a component that provides showToast",
    );
  }

  return showToast;
};
