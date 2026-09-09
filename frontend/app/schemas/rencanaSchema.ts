import { z } from 'zod'

export const rencanaSchema = z.object({
  nama_customer: z
    .string()
    .trim()
    .min(1, 'Nama nasabah wajib diisi')
    .min(3, 'Nama nasabah minimal 3 karakter'),
  hp: z
    .string()
    .trim()
    .min(1, 'Nomor HP / WhatsApp wajib diisi')
    .min(9, 'Nomor HP minimal 9 digit')
    .regex(/^[0-9+\-\s]+$/, 'Nomor HP hanya boleh berisi angka'),
  status: z
    .string()
    .min(1, 'Status pipeline wajib dipilih'),
  ket_rencana: z
    .string()
    .optional()
    .nullable(),
  ket_aktivitas: z
    .string()
    .optional()
    .nullable(),
  id_sumbercust: z
    .number()
    .optional()
    .nullable()
})

export type RencanaSchemaType = z.infer<typeof rencanaSchema>

export const validateRencana = (data: unknown) => {
  return rencanaSchema.safeParse(data)
}
