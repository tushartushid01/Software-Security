import { yupResolver } from "@hookform/resolvers/yup/src/yup.js";
import {
  Button,
  Dialog,
  DialogTitle,
  Stack,
  TextField,
  Typography,
} from "@mui/material";
import { useCallback, useState } from "react";
import { useForm } from "react-hook-form";
import * as Yup from "yup";

const schema = Yup.object().shape({
  feedback: Yup.string()
    .min(10, "atleast 10 character required")
    .required("feed required"),
});

export const FeedbackDialog = ({ open, handleClose }: any) => {
  const token = localStorage.getItem("token");
  const [errorMessage, setErrorMessage] = useState<string>("");
  const {
    register,
    handleSubmit,
    watch,
    formState: { errors },
    reset
  } = useForm<any>({
    defaultValues: {
      feedback: "",
    },
    mode: "all",
    resolver: yupResolver(schema),
  });

  const handleFeedbackSubmit = useCallback(async () => {
    try {
      const response = await fetch(`http://localhost:8080/v1/auth/feedback`, {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
          Authorization: `${token}`,
        },
        body: JSON.stringify({
          description: watch("feedback"),
        }),
      });
      if (response.status === 200) {
        setErrorMessage("")
        handleClose();
      }else{
        setErrorMessage("Something went wrong");
      }
    } catch (error) {
      console.log(error);
      setErrorMessage("Something went wrong");
    }
  }, [handleClose, token, watch]);

  return (
    <Dialog open={open} onClose={()=>{handleClose()
      reset()
      setErrorMessage("")
    }} maxWidth={"sm"}>
      <DialogTitle textAlign="center">Feedback</DialogTitle>
      <Stack
        component="form"
        onSubmit={handleSubmit(handleFeedbackSubmit)}
        spacing={2}
        mt={2}
        width={400}
        p={4}
      >
        <TextField
          fullWidth
          label="Feedback"
          value={watch("feedback")}
          {...register("feedback")}
          helperText={errors && (errors.feedback?.message as string)}
          sx={{ "& .MuiFormHelperText-root": { color: "red" } }}
        />
        {errorMessage && <Typography color="error">{errorMessage}</Typography>}
        <Button variant="contained" type="submit">
          Submit
        </Button>
      </Stack>
    </Dialog>
  );
};
