import { yupResolver } from "@hookform/resolvers/yup/src/yup.js";
import {
  Button,
  Dialog,
  DialogContent,
  DialogTitle,
  Stack,
  TextField,
  Typography,
} from "@mui/material";
import { useCallback, useState } from "react";
import { useForm } from "react-hook-form";
import * as Yup from "yup";

const schema = Yup.object().shape({
  title: Yup.string().required("Title is required"),
  description: Yup.string().required("Description is required"),
  price: Yup.string().required("Price is required"),
});

export const SellingDialog = ({ open, handleClose , getProductDetails}: any) => {
  const [errorMessage, setErrorMessage] = useState<string>("");
  const token = localStorage.getItem("token");
  const {
    register,
    handleSubmit,
    watch,
    formState: { errors },
    reset
  } = useForm<any>({
    defaultValues: {
      title: "",
      description: "",
      price: "",
    },
    mode: "all",
    resolver: yupResolver(schema),
  });

  const handleSubmitSellingForm = useCallback(async () => {
    try {
      const response = await fetch(
        `http://localhost:8080/v1/auth/admin/sell-product`,
        {
          method: "POST",
          headers: {
            "Content-Type": "application/json",
            Authorization: `${token}`,
          },
          body: JSON.stringify({
            name: watch("title"),
            description: watch("description"),
            price: Number(watch("price")),
          }),
        }
      );
      if (response.status === 200) {
        reset()
        handleClose();
        getProductDetails()
        setErrorMessage("")
      }else{
        setErrorMessage("Something went wrong");
      }
    } catch (error) {
      console.log(error);
      
    }
  }, [handleClose, token, watch, getProductDetails]);

  return (
    <Dialog open={open} onClose={()=>{handleClose()
      reset()
      setErrorMessage("")
    }} fullWidth maxWidth={"sm"}>
      <DialogTitle textAlign="center">Selling Form</DialogTitle>
      <DialogContent>
        <Stack
          component={"form"}
          onSubmit={handleSubmit(handleSubmitSellingForm)}
          spacing={2}
          mt={2}
          width={"100%"}
        >
          <TextField
            fullWidth
            label="Product Title"
            value={watch("title")}
            {...register("title")}
            helperText={errors && (errors.title?.message as string)}
            sx={{ "& .MuiFormHelperText-root": { color: "red" } }}
          />
          <TextField
            fullWidth
            label="Product Description"
            value={watch("description")}
            {...register("description")}
            helperText={errors && (errors.description?.message as string)}
            sx={{ "& .MuiFormHelperText-root": { color: "red" } }}
          />
          <TextField
            fullWidth
            label="Price"
            value={watch("price")}
            {...register("price")}
            helperText={errors && (errors.price?.message as string)}
            sx={{ "& .MuiFormHelperText-root": { color: "red" } }}
          />
           {errorMessage && <Typography color="error">{errorMessage}</Typography>}
          <Button variant="contained" type="submit">
            Submit
          </Button>
        </Stack>
      </DialogContent>
    </Dialog>
  );
};
